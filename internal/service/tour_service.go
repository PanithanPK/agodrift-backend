package service

import (
	"agodrift/internal/config"
	"agodrift/internal/model"
	"agodrift/internal/repository"
)

// TourService provides business logic for tours.
type TourService struct {
	repo repository.TourRepository
}

func NewTourService() *TourService {
	return &TourService{
		repo: repository.NewMySQLTourRepo(config.GetDB()),
	}
}

func (s *TourService) List() []model.Tour {
	return s.repo.List()
}

func (s *TourService) Get(id int) (model.Tour, bool) {
	return s.repo.Get(id)
}

func (s *TourService) Create(t model.Tour) model.Tour {
	return s.repo.Create(t)
}
