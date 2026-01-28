package service

import (
	"time"

	"agodrift/internal/config"
	"agodrift/internal/model"
	"agodrift/internal/repository"
)

type BookingService struct {
	repo repository.BookingRepository
}

func NewBookingService() *BookingService {
	return &BookingService{repo: repository.NewMySQLBookingRepo(config.GetDB())}
}

func (s *BookingService) Create(userID int, hotelID int, checkIn time.Time, checkOut time.Time, adults int, children int, rooms int) (model.Booking, error) {
	return s.repo.Create(userID, hotelID, checkIn, checkOut, adults, children, rooms)
}

func (s *BookingService) ListByUserID(userID int) ([]model.Booking, error) {
	return s.repo.ListByUserID(userID)
}
