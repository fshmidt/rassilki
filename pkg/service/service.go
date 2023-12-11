package service

import (
	"github.com/fshmidt/rassilki"
	"github.com/fshmidt/rassilki/pkg/repository"
)

type Client interface {
	Create(client rassilki.Client) (id int, err error)
	Get(id int) (client rassilki.Client, err error)
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
	GetById(id int) (rassilki.Rassilka, error)
}

type Messages interface {
	CreateTable(id int, clientsId []rassilki.Client) error
	DropTable(id int) error
	GetClientsList(id int) ([]rassilki.Client, error)
	GetMessageText(id int) (string, error)
	GetMessage(subscriberID, rassilkaId int) (rassilki.Message, error)
	IsStarted(id int) (bool, error)
	RenewTable(input rassilki.UpdateRassilka, id int) error
	SendMessageToSubscriber(msgID, subscriberID, rassilkaId int, text string) error
	GetRassilkiReview(ids []int) ([]rassilki.RassilkaReview, error)
	GetRassilkaReviewById(id int) (rassilki.RassilkaReview, error)
}

type Service struct {
	Client
	Rassilka
	Messages
}

func NewService(repos *repository.Repository, jwt string) *Service {
	return &Service{
		Client:   NewClientService(repos.Client),
		Rassilka: NewRassilkaService(repos.Rassilka),
		Messages: NewMessagesService(repos.Messages, jwt),
	}
}
