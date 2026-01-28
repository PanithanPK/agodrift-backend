package model

// User represents a simple user record
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`    // plaintext for demo; in real app store hashed password
	Role     string `json:"role"` // "admin" or "user"
}
