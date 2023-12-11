package service

import (
	"github.com/fshmidt/rassilki"
	"github.com/fshmidt/rassilki/pkg/repository"
)

type RassilkaService struct {
	repo repository.Rassilka
}

func NewRassilkaService(repo repository.Rassilka) *RassilkaService {

	return &RassilkaService{repo: repo}
}

func (s *RassilkaService) Create(rassilka rassilki.Rassilka) (int, error) {

	return s.repo.Create(rassilka)
}

func (s *RassilkaService) Delete(id int) error {

	return s.repo.Delete(id)
}

func (s *RassilkaService) Update(input rassilki.UpdateRassilka, id int) error {

	return s.repo.Update(input, id)
}

func (s *RassilkaService) GetById(id int) (Rassilka rassilki.Rassilka, err error) {

	return s.repo.GetById(id)
}

func (s *RassilkaService) CheckActive() ([]int, error) {

	return s.repo.CheckActive()
}

func (s *RassilkaService) CheckUpdated() ([]int, error) {

	return s.repo.CheckUpdated()
}

func (s *RassilkaService) CheckRecreated() ([]int, error) {

	return s.repo.CheckRecreated()
}
func (s *RassilkaService) GetAll() ([]int, error) {

	return s.repo.GetAll()
}
