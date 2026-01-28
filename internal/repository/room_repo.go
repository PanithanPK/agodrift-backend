package repository

import (
	"database/sql"
	"sync"

	"agodrift/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

type RoomRepository interface {
	List() []model.Room
	Get(id int) (model.Room, bool)
	Create(r model.Room) model.Room
}

type inMemoryRoomRepo struct {
	mu    sync.RWMutex
	rooms map[int]model.Room
	next  int
}

func NewInMemoryRoomRepo() *inMemoryRoomRepo {
	r := &inMemoryRoomRepo{
		rooms: make(map[int]model.Room),
		next:  1,
	}
	r.Create(model.Room{Name: "Demo Hotel", Description: "Demo", Location: "Demo", Destination: "Demo", Rating: 4.5, Reviews: 10, PriceCents: 15000, Amenities: "Wi-Fi", Featured: true, MaxAdults: 2, MaxChildren: 1, RoomsTotal: 10, RoomsAvailable: 5, Status: "active"})
	return r
}

func (r *inMemoryRoomRepo) List() []model.Room {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Room, 0, len(r.rooms))
	for _, rm := range r.rooms {
		out = append(out, rm)
	}
	return out
}

func (r *inMemoryRoomRepo) Get(id int) (model.Room, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rm, ok := r.rooms[id]
	return rm, ok
}

func (r *inMemoryRoomRepo) Create(rm model.Room) model.Room {
	r.mu.Lock()
	defer r.mu.Unlock()
	rm.ID = r.next
	r.next++
	r.rooms[rm.ID] = rm
	return rm
}

type mysqlRoomRepo struct {
	db *sql.DB
}

func NewMySQLRoomRepo(db *sql.DB) *mysqlRoomRepo {
	return &mysqlRoomRepo{db: db}
}

func (r *mysqlRoomRepo) List() []model.Room {
	rows, err := r.db.Query("SELECT id, name, description, location, destination, rating, reviews, price_cents, original_price_cents, amenities, featured, max_adults, max_children, rooms_total, rooms_available, status FROM hotels")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var rm model.Room
		var original sql.NullInt64
		var featuredInt int
		if err := rows.Scan(
			&rm.ID,
			&rm.Name,
			&rm.Description,
			&rm.Location,
			&rm.Destination,
			&rm.Rating,
			&rm.Reviews,
			&rm.PriceCents,
			&original,
			&rm.Amenities,
			&featuredInt,
			&rm.MaxAdults,
			&rm.MaxChildren,
			&rm.RoomsTotal,
			&rm.RoomsAvailable,
			&rm.Status,
		); err != nil {
			continue
		}
		rm.Featured = featuredInt == 1
		if original.Valid {
			v := int(original.Int64)
			rm.OriginalPriceCents = &v
		}
		rooms = append(rooms, rm)
	}
	return rooms
}

func (r *mysqlRoomRepo) Get(id int) (model.Room, bool) {
	var rm model.Room
	var original sql.NullInt64
	var featuredInt int
	err := r.db.QueryRow("SELECT id, name, description, location, destination, rating, reviews, price_cents, original_price_cents, amenities, featured, max_adults, max_children, rooms_total, rooms_available, status FROM hotels WHERE id = ?", id).Scan(
		&rm.ID,
		&rm.Name,
		&rm.Description,
		&rm.Location,
		&rm.Destination,
		&rm.Rating,
		&rm.Reviews,
		&rm.PriceCents,
		&original,
		&rm.Amenities,
		&featuredInt,
		&rm.MaxAdults,
		&rm.MaxChildren,
		&rm.RoomsTotal,
		&rm.RoomsAvailable,
		&rm.Status,
	)
	if err != nil {
		return rm, false
	}
	rm.Featured = featuredInt == 1
	if original.Valid {
		v := int(original.Int64)
		rm.OriginalPriceCents = &v
	}
	return rm, true
}

func (r *mysqlRoomRepo) Create(rm model.Room) model.Room {
	original := sql.NullInt64{}
	if rm.OriginalPriceCents != nil {
		original = sql.NullInt64{Int64: int64(*rm.OriginalPriceCents), Valid: true}
	}
	featured := 0
	if rm.Featured {
		featured = 1
	}
	if rm.Status == "" {
		rm.Status = "active"
	}
	if rm.RoomsTotal == 0 {
		rm.RoomsTotal = 1
	}
	if rm.RoomsAvailable == 0 {
		rm.RoomsAvailable = rm.RoomsTotal
	}
	if rm.MaxAdults == 0 {
		rm.MaxAdults = 1
	}
	result, err := r.db.Exec(
		"INSERT INTO hotels (name, description, location, destination, rating, reviews, price_cents, original_price_cents, amenities, featured, max_adults, max_children, rooms_total, rooms_available, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rm.Name,
		rm.Description,
		rm.Location,
		rm.Destination,
		rm.Rating,
		rm.Reviews,
		rm.PriceCents,
		original,
		rm.Amenities,
		featured,
		rm.MaxAdults,
		rm.MaxChildren,
		rm.RoomsTotal,
		rm.RoomsAvailable,
		rm.Status,
	)
	if err != nil {
		return rm
	}
	id, _ := result.LastInsertId()
	rm.ID = int(id)
	return rm
}
