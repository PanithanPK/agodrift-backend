package tests

import (
	"testing"

	"agodrift/internal/repository"
	"agodrift/internal/service"
)

func TestListRooms(t *testing.T) {
	s := service.NewRoomServiceWithRepo(repository.NewInMemoryRoomRepo())
	list := s.List()
	if len(list) < 1 {
		t.Fatalf("expected seeded rooms, got %d", len(list))
	}
}
