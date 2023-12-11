package service

import (
	"github.com/fshmidt/rassilki"
	"github.com/fshmidt/rassilki/pkg/repository"
)

type ClientService struct {
	repo repository.Client
}

func NewClientService(repo repository.Client) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) Create(client rassilki.Client) (int, error) {
	return s.repo.Create(client)
}

func (s *ClientService) Get(id int) (client rassilki.Client, err error) {
	return s.repo.Get(id)
}

func (s *ClientService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ClientService) Update(input rassilki.UpdateClient, id int) error {
	return s.repo.Update(input, id)
}
