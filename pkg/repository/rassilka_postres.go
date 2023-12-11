package repository

import (
	"fmt"
	"github.com/fshmidt/rassilki"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"time"
)

const (
	time_changes = iota
	recreate_messages_table
	add_clients_to_messages_table
	no_changes
)

type RassilkaPostgres struct {
	db *sqlx.DB
}

func NewRassilkaPostgres(db *sqlx.DB) *RassilkaPostgres {
	return &RassilkaPostgres{db: db}
}

func (r *RassilkaPostgres) Create(rassilka rassilki.Rassilka) (int, error) {

	var id int
	query := fmt.Sprintf("INSERT INTO %s (start_time , message, filter, end_time, supplemented, recreated) values($1,$2,$3,$4,$5,$6) RETURNING id", rassilkiTable)

	filterArray := pq.Array(rassilka.Filter)
	row := r.db.QueryRow(query, rassilka.StartTime, rassilka.Message, filterArray, rassilka.EndTime, false, false)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RassilkaPostgres) Delete(id int) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", rassilkiTable)
	_, err = r.db.Exec(query, id)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		DROP TABLE messages_%d `, id)

	_, err = r.db.Exec(query)
	return err
}

func (r *RassilkaPostgres) Update(input rassilki.UpdateRassilka, id int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.StartTime != nil {
		setValues = append(setValues, fmt.Sprintf("start_time=$%d", argId))
		args = append(args, *input.StartTime)
		argId++
	}

	if input.Message != nil {
		setValues = append(setValues, fmt.Sprintf("message=$%d", argId))
		args = append(args, *input.Message)
		argId++
	}

	if input.EndTime != nil {
		setValues = append(setValues, fmt.Sprintf("end_time=$%d", argId))
		args = append(args, *input.EndTime)
		argId++
	}

	if *input.Supplemented == true {
		setValues = append(setValues, fmt.Sprintf("supplemented=$%d", argId))
		args = append(args, *input.Supplemented)
		argId++
	} else if *input.Recreated == true {
		setValues = append(setValues, fmt.Sprintf("recreated=$%d", argId))
		args = append(args, *input.Recreated)
		argId++
	}

	if input.Filter != nil {
		setValues = append(setValues, fmt.Sprintf("filter=$%d", argId))
		filterArray := pq.Array(*input.Filter)
		args = append(args, filterArray)
	} else {
		setValues = append(setValues, fmt.Sprintf("filter=$%d", argId))
		args = append(args, pq.Array([]string{}))
	}
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d",
		rassilkiTable, setQuery, argId)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *RassilkaPostgres) CheckActive() ([]int, error) {
	currentTime := time.Now()
	query := fmt.Sprintf("SELECT id FROM %s WHERE start_time <= $1 AND end_time > $2", rassilkiTable)
	rows, err := r.db.Query(query, currentTime, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		activeIDs = append(activeIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activeIDs, nil
}

func (r *RassilkaPostgres) CheckRecreated() ([]int, error) {

	query := fmt.Sprintf("SELECT id FROM %s WHERE recreated = $1", rassilkiTable)
	rows, err := r.db.Query(query, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recreatedIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		recreatedIDs = append(recreatedIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	//reset recreated

	if len(recreatedIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(recreatedIDs))
	for i := range recreatedIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query = fmt.Sprintf("UPDATE %s SET recreated=false WHERE id = ANY($1)", rassilkiTable)

	_, err = r.db.Exec(query, pq.Array(recreatedIDs))
	if err != nil {
		return nil, err
	}

	return recreatedIDs, nil
}
func (r *RassilkaPostgres) CheckUpdated() ([]int, error) {

	query := fmt.Sprintf("SELECT id FROM %s WHERE supplemented = true", rassilkiTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var updatedIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		updatedIDs = append(updatedIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	//reset updated
	if len(updatedIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(updatedIDs))
	for i := range updatedIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query = fmt.Sprintf("UPDATE %s SET supplemented=false WHERE id = ANY($1)", rassilkiTable)

	_, err = r.db.Exec(query, pq.Array(updatedIDs))
	if err != nil {
		return nil, err
	}
	return updatedIDs, nil
}

func (r *RassilkaPostgres) GetAll() ([]int, error) {

	query := fmt.Sprintf("SELECT id FROM %s", rassilkiTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		//logrus.Println("Error iterating over rows: %v", err)
		return nil, err
	}
	return ids, nil
}

func (r *RassilkaPostgres) GetById(id int) (rassilka rassilki.Rassilka, err error) {

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", rassilkiTable)
	err = r.db.Get(&rassilka, query, id)

	return rassilka, err
}
