package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbOnce sync.Once
	dbInst *sql.DB
)

// Get reads an env var with fallback default.
func Get(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// GetDB returns a database connection.
func GetDB() *sql.DB {
	dbOnce.Do(func() {
		dsn := strings.TrimSpace(Get("MYSQL_DSN", ""))
		if dsn == "" {
			dsn = strings.TrimSpace(Get("DATABASE_URL", ""))
		}
		if dsn == "" {
			user := Get("DB_USER", "user")
			pass := Get("DB_PASS", "123456")
			host := Get("DB_HOST", "127.0.0.1")
			port := Get("DB_PORT", "3306")
			name := Get("DB_NAME", "agodrift")
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", user, pass, host, port, name)
		} else if strings.HasPrefix(strings.ToLower(dsn), "mysql://") {
			// Accept mysql://user:pass@host:3306/dbname style and convert to go-sql-driver/mysql DSN.
			u, err := url.Parse(dsn)
			if err != nil {
				log.Fatal(err)
			}
			user := ""
			pass := ""
			if u.User != nil {
				user = u.User.Username()
				pass, _ = u.User.Password()
			}
			host := u.Hostname()
			port := u.Port()
			if port == "" {
				port = "3306"
			}
			name := strings.TrimPrefix(u.Path, "/")
			q := u.Query()
			if q.Get("parseTime") == "" {
				q.Set("parseTime", "true")
			}
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pass, host, port, name, q.Encode())
		}

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}
		// Retry connection until successful
		for i := 0; i < 30; i++ {
			if err := db.Ping(); err == nil {
				break
			}
			log.Println("Waiting for database connection...")
			time.Sleep(2 * time.Second)
		}
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		dbInst = db
	})

	return dbInst
}
