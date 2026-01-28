package service

import (
	"agodrift/internal/config"
	"agodrift/internal/model"
	"agodrift/internal/repository"
)

type RoomService struct {
	repo repository.RoomRepository
}

func NewRoomService() *RoomService {
	return &RoomService{repo: repository.NewMySQLRoomRepo(config.GetDB())}
}

func NewRoomServiceWithRepo(repo repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) List() []model.Room {
	return s.repo.List()
}

func (s *RoomService) Get(id int) (model.Room, bool) {
	return s.repo.Get(id)
}

func (s *RoomService) Create(r model.Room) model.Room {
	return s.repo.Create(r)
}
