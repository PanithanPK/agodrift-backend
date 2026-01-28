package repository

import (
	"database/sql"
	"sync"

	"agodrift/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

// TourRepository defines methods for tour storage.
type TourRepository interface {
	List() []model.Tour
	Get(id int) (model.Tour, bool)
	Create(t model.Tour) model.Tour
}

type inMemoryTourRepo struct {
	mu    sync.RWMutex
	tours map[int]model.Tour
	next  int
}

func NewInMemoryTourRepo() *inMemoryTourRepo {
	r := &inMemoryTourRepo{
		tours: make(map[int]model.Tour),
		next:  1,
	}
	// seed with sample data
	r.Create(model.Tour{Title: "Bangkok Highlights", Description: "City tour", PriceCents: 15000})
	r.Create(model.Tour{Title: "Chiang Mai Trek", Description: "Jungle trek", PriceCents: 25000})
	return r
}

func (r *inMemoryTourRepo) List() []model.Tour {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Tour, 0, len(r.tours))
	for _, t := range r.tours {
		out = append(out, t)
	}
	return out
}

func (r *inMemoryTourRepo) Get(id int) (model.Tour, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tours[id]
	return t, ok
}

func (r *inMemoryTourRepo) Create(t model.Tour) model.Tour {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = r.next
	r.next++
	r.tours[t.ID] = t
	return t
}

type mysqlTourRepo struct {
	db *sql.DB
}

func NewMySQLTourRepo(db *sql.DB) *mysqlTourRepo {
	return &mysqlTourRepo{db: db}
}

func (r *mysqlTourRepo) List() []model.Tour {
	rows, err := r.db.Query("SELECT id, title, description, price_cents FROM tours")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var tours []model.Tour
	for rows.Next() {
		var t model.Tour
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.PriceCents); err != nil {
			continue
		}
		tours = append(tours, t)
	}
	return tours
}

func (r *mysqlTourRepo) Get(id int) (model.Tour, bool) {
	var t model.Tour
	err := r.db.QueryRow("SELECT id, title, description, price_cents FROM tours WHERE id = ?", id).Scan(&t.ID, &t.Title, &t.Description, &t.PriceCents)
	if err != nil {
		return t, false
	}
	return t, true
}

func (r *mysqlTourRepo) Create(t model.Tour) model.Tour {
	result, err := r.db.Exec("INSERT INTO tours (title, description, price_cents) VALUES (?, ?, ?)", t.Title, t.Description, t.PriceCents)
	if err != nil {
		return t
	}
	id, _ := result.LastInsertId()
	t.ID = int(id)
	return t
}
