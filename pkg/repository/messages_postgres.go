package repository

import (
	"errors"
	"fmt"
	"github.com/fshmidt/rassilki"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type MessagesPostgres struct {
	db *sqlx.DB
}

func NewMessagesPostgres(db *sqlx.DB) *MessagesPostgres {
	return &MessagesPostgres{db: db}
}

func (s *MessagesPostgres) CreateTable(id int, clients []rassilki.Client) error {
	started, _ := s.IsStarted(id)
	if started {
		return errors.New("THIS SHIT IS STARTED ALREADY")
	}
	tx, err := s.db.Begin()
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

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS messages_%d (
		id INT PRIMARY KEY DEFAULT nextval('messages_id_seq'),
		sent_at TIMESTAMPTZ,
		status BOOLEAN NOT NULL,
		ras_id INT NOT NULL,
		client_id INT NOT NULL
	)`, id)

	_, err = s.db.Exec(query)
	if err != nil {
		return err
	}

	subquery := fmt.Sprintf("SELECT id FROM %s WHERE id = ANY($1::int[])", clientsTable)

	var clientsId []int
	for _, client := range clients {
		clientsId = append(clientsId, client.Id)
	}

	rows, err := s.db.Query(subquery, pq.Array(clientsId))
	if err != nil {
		return err
	}
	defer rows.Close()

	var foundIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		foundIDs = append(foundIDs, id)
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO messages_%d (sent_at, status, ras_id, client_id)
		SELECT NULL, false, $1, id
		FROM %s
		WHERE id = ANY($2::int[])`, id, clientsTable)
	_, err = tx.Exec(insertQuery, id, pq.Array(foundIDs))
	if err != nil {
		return err
	}
	return nil
}

func (s *MessagesPostgres) DropTable(id int) error {
	query := fmt.Sprintf(`
		DROP TABLE messages_%d `, id)

	_, err := s.db.Exec(query)
	return err
}

func (s *MessagesPostgres) GetClientsList(rassilkaId int) ([]rassilki.Client, error) {
	var clients []rassilki.Client

	var rassilka rassilki.Rassilka
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, rassilkiTable)
	if err := s.db.Get(&rassilka, query, rassilkaId); err != nil {
		return nil, err
	}

	for _, tag := range rassilka.Filter {
		subQuery := fmt.Sprintf("SELECT * FROM %s WHERE tag = $1", clientsTable)
		rows, err := s.db.Query(subQuery, tag)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var client rassilki.Client
			if err := rows.Scan(&client.Id, &client.Phone, &client.Code, &client.Tag, &client.Timezone); err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			clients = append(clients, client)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}
	return clients, nil
}

func (s *MessagesPostgres) GetMessageText(rassilkaId int) (string, error) {

	var message string
	query := fmt.Sprintf("SELECT message FROM %s WHERE id=$1", rassilkiTable)
	err := s.db.Get(&message, query, rassilkaId)

	return message, err
}

func (s *MessagesPostgres) GetMessage(clientId, rassilkaId int) (rassilki.Message, error) {

	var message rassilki.Message
	query := fmt.Sprintf("SELECT * FROM %s WHERE client_id=$1", "messages_"+strconv.Itoa(rassilkaId))
	err := s.db.Get(&message, query, clientId)

	return message, err
}

func (s *MessagesPostgres) IsStarted(id int) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = $1
        );
    `
	err := s.db.QueryRow(query, "messages_"+strconv.Itoa(id)).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *MessagesPostgres) UpdateMessageStatus(messageId, rassilkaId int) error {

	currentTime := time.Now()
	query := fmt.Sprintf("UPDATE %s SET status=TRUE, sent_at=$1 WHERE id = $2", "messages_"+strconv.Itoa(rassilkaId))

	_, err := s.db.Exec(query, currentTime, messageId)
	return err
}

func (s *MessagesPostgres) RenewTable(input rassilki.UpdateRassilka, id int) error {

	query := fmt.Sprintf("SELECT id FROM %s WHERE tag = ANY($1::text[])", clientsTable)
	rows, err := s.db.Query(query, pq.Array(*input.Filter))
	if err != nil {
		return err
	}
	defer rows.Close()

	var foundIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		foundIDs = append(foundIDs, id)
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO messages_%d (sent_at, status, ras_id, client_id)
		SELECT NULL, false, $1, id
		FROM %s
		WHERE id = ANY($2::int[])`, id, clientsTable)
	_, err = s.db.Exec(insertQuery, id, pq.Array(foundIDs))
	if err != nil {
		return err
	}
	return nil
}

func (s *MessagesPostgres) GetRassilkiReview(ids []int) ([]rassilki.RassilkaReview, error) {

	var reviews []rassilki.RassilkaReview

	for _, id := range ids {

		started, _ := s.IsStarted(id)
		if !started {
			continue
		}
		tableName := fmt.Sprintf("messages_%d", id)

		var total int
		err := s.db.Get(&total, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
		if err != nil {
			return nil, err
		}

		var sent int
		err = s.db.Get(&sent, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE status = true", tableName))
		if err != nil {
			return nil, err
		}

		notSent := total - sent

		review := rassilki.RassilkaReview{
			Id:      id,
			Total:   total,
			Sent:    sent,
			NotSent: notSent,
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (s *MessagesPostgres) GetRassilkaReviewById(id int) (rassilki.RassilkaReview, error) {

	review := rassilki.RassilkaReview{
		Id:      id,
		Total:   0,
		Sent:    0,
		NotSent: 0,
	}
	started, _ := s.IsStarted(id)
	if !started {
		return review, errors.New("rassilka hasn't started or created yet")
	}

	tableName := fmt.Sprintf("messages_%d", id)

	var total int
	err := s.db.Get(&total, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
	if err != nil {
		return review, err
	}

	var sent int
	err = s.db.Get(&sent, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE status = true", tableName))
	if err != nil {
		return review, err
	}

	notSent := total - sent

	review.Total, review.Sent, review.NotSent = total, sent, notSent

	return review, nil
}
