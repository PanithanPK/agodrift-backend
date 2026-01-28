package tests

import (
	"testing"

	"agodrift/internal/service"
)

func TestListTours(t *testing.T) {
	s := service.NewTourService()
	list := s.List()
	if len(list) < 1 {
		t.Fatalf("expected seeded tours, got %d", len(list))
	}
}
