package repository

import (
	"database/sql"
	"sync"

	"agodrift/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

// UserRepository defines methods for user storage.
type UserRepository interface {
	GetByEmail(email string) (model.User, bool)
	Create(u model.User) model.User
}

type inMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]model.User // keyed by email
	next  int
}

// NewInMemoryUserRepo returns a seeded in-memory repo with a default admin and user.
func NewInMemoryUserRepo() *inMemoryUserRepo {
	r := &inMemoryUserRepo{
		users: make(map[string]model.User),
		next:  1,
	}
	// seed users (passwords are plain for demo only)
	r.Create(model.User{Name: "Admin User", Email: "admin@agodrift.dev", Password: "adminpass", Role: "admin"})
	r.Create(model.User{Name: "Alice Traveler", Email: "alice@example.com", Password: "userpass", Role: "user"})
	return r
}

func (r *inMemoryUserRepo) GetByEmail(email string) (model.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[email]
	return u, ok
}

func (r *inMemoryUserRepo) Create(u model.User) model.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.next
	r.next++
	r.users[u.Email] = u
	return u
}

type mysqlUserRepo struct {
	db *sql.DB
}

func NewMySQLUserRepo(db *sql.DB) *mysqlUserRepo {
	return &mysqlUserRepo{db: db}
}

func (r *mysqlUserRepo) GetByEmail(email string) (model.User, bool) {
	var u model.User
	err := r.db.QueryRow("SELECT id, name, email, password, role FROM users WHERE email = ?", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return u, false
	}
	return u, true
}

func (r *mysqlUserRepo) Create(u model.User) model.User {
	result, err := r.db.Exec("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)", u.Name, u.Email, u.Password, u.Role)
	if err != nil {
		return u
	}
	id, _ := result.LastInsertId()
	u.ID = int(id)
	return u
}
