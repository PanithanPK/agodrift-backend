package repository

import (
	"database/sql"
	"sync"

	"agodrift/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

// UserRepository defines methods for user storage.
type UserRepository interface {
	GetByUsername(username string) (model.User, bool)
	Create(u model.User) model.User
}

type inMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]model.User // keyed by username
	next  int
}

// NewInMemoryUserRepo returns a seeded in-memory repo with a default admin and user.
func NewInMemoryUserRepo() *inMemoryUserRepo {
	r := &inMemoryUserRepo{
		users: make(map[string]model.User),
		next:  1,
	}
	// seed users (passwords are plain for demo only)
	r.Create(model.User{Username: "admin", Password: "adminpass", Role: "admin"})
	r.Create(model.User{Username: "alice", Password: "userpass", Role: "user"})
	return r
}

func (r *inMemoryUserRepo) GetByUsername(username string) (model.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[username]
	return u, ok
}

func (r *inMemoryUserRepo) Create(u model.User) model.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.next
	r.next++
	r.users[u.Username] = u
	return u
}

type mysqlUserRepo struct {
	db *sql.DB
}

func NewMySQLUserRepo(db *sql.DB) *mysqlUserRepo {
	return &mysqlUserRepo{db: db}
}

func (r *mysqlUserRepo) GetByUsername(username string) (model.User, bool) {
	var u model.User
	err := r.db.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		return u, false
	}
	return u, true
}

func (r *mysqlUserRepo) Create(u model.User) model.User {
	result, err := r.db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", u.Username, u.Password, u.Role)
	if err != nil {
		return u
	}
	id, _ := result.LastInsertId()
	u.ID = int(id)
	return u
}
