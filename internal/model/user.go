package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`    // plaintext for demo; in real app store hashed password
	Role     string `json:"role"` // "admin" or "user"
}
