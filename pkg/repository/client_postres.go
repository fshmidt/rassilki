package repository

import (
	"fmt"
	"github.com/fshmidt/rassilki"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type ClientPostgres struct {
	db *sqlx.DB
}

func NewClientPostgres(db *sqlx.DB) *ClientPostgres {
	return &ClientPostgres{db: db}
}

func (r *ClientPostgres) Create(client rassilki.Client) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var clientId int
	query := fmt.Sprintf("INSERT INTO %s (phone, code, tag, timezone) values($1,$2,$3,$4) RETURNING id", clientsTable)

	row := r.db.QueryRow(query, client.Phone, client.Code, client.Tag, client.Timezone)
	if err := row.Scan(&clientId); err != nil {
		return 0, err
	}

	subQuery := fmt.Sprintf("SELECT id FROM %s WHERE $1 = ANY(filter)", rassilkiTable)
	rows, err := r.db.Query(subQuery, client.Tag)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var rasIds []int
	for rows.Next() {
		var rasId int
		if err := rows.Scan(&rasId); err != nil {
			return 0, err
		}

		rasIds = append(rasIds, rasId)
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	for _, rasId := range rasIds {
		msgQuery := fmt.Sprintf("INSERT INTO %s (status, ras_id, client_id) values($1,$2, $3)", "messages_"+strconv.Itoa(rasId))

		_, err := r.db.Exec(msgQuery, false, rasId, clientId)
		if err != nil {
			return 0, err
		}
	}

	return clientId, nil
}

func (r *ClientPostgres) Get(id int) (rassilki.Client, error) {
	var client rassilki.Client
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", clientsTable)
	err := r.db.Get(&client, query, id)

	return client, err
}

func (r *ClientPostgres) Delete(id int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", clientsTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *ClientPostgres) Update(input rassilki.UpdateClient, id int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *input.Phone)
		argId++
	}

	if input.Tag != nil {
		setValues = append(setValues, fmt.Sprintf("tag=$%d", argId))
		args = append(args, *input.Tag)
		argId++
	}

	if input.Code != nil {
		setValues = append(setValues, fmt.Sprintf("code=$%d", argId))
		args = append(args, *input.Code)
		argId++
	}

	if input.Timezone != nil {
		setValues = append(setValues, fmt.Sprintf("timezone=$%d", argId))
		args = append(args, *input.Timezone)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d",
		clientsTable, setQuery, argId)
	args = append(args, id)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
