package repository

import (
	"github.com/fshmidt/rassilki"
	"github.com/jmoiron/sqlx"
)

type Client interface {
	Create(client rassilki.Client) (int, error)
	Get(id int) (rassilki.Client, error)
	Delete(id int) error
	Update(input rassilki.UpdateClient, id int) error
}

type Rassilka interface {
	Create(Rassilka rassilki.Rassilka) (int, error)
	Update(input rassilki.UpdateRassilka, id int) error
	Delete(id int) error
	CheckActive() ([]int, error)
	CheckRecreated() ([]int, error)
	CheckUpdated() ([]int, error)
	GetAll() ([]int, error)
	GetById(id int) (Rassilka rassilki.Rassilka, err error)
}

type Messages interface {
	CreateTable(id int, clients []rassilki.Client) error
	DropTable(id int) error
	GetClientsList(rassilkaId int) ([]rassilki.Client, error)
	GetMessageText(rassilkaId int) (string, error)
	GetMessage(subscriberID, rassilkaId int) (rassilki.Message, error)
	UpdateMessageStatus(messageId, rassilkaId int) error
	IsStarted(id int) (bool, error)
	RenewTable(input rassilki.UpdateRassilka, id int) error
	GetRassilkiReview(ids []int) ([]rassilki.RassilkaReview, error)
	GetRassilkaReviewById(id int) (rassilki.RassilkaReview, error)
}

type Repository struct {
	Client
	Rassilka
	Messages
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Client:   NewClientPostgres(db),
		Rassilka: NewRassilkaPostgres(db),
		Messages: NewMessagesPostgres(db),
	}
}
