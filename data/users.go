package data

import (
	"database/sql"
	"errors"
	"time"

	"fmt"

	"github.com/google/uuid"
)

var ErrInvalidEmailOrPassword = errors.New("users: Invalid email or password")
var ErrEmailNotPresent = errors.New("users: Email not registered")

// NewUser is the type received from mobile apps before being saved into the db
type NewUser struct {
	Name     string
	Email    string
	Password string
}

// User contains info on our user
type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	Password       []byte
	EmailConfirmed bool
	Created        time.Time
}

// UserStore is our interface defining methods for concrete service
type UserStore interface {
	GetUser(id string) (User, error)
	ConfEmail(id string) error
}

// UserService is our psql implementation of UserStore
type UserService struct {
	db *sql.DB
}

// InitUserService initializes a new UserService
func InitUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// GetUser returns a user from the database
func (us *UserService) GetUser(id string) (User, error) {
	var email sql.NullString
	var name sql.NullString

	row := us.db.QueryRow("SELECT user_id, email, name, password, email_confirmed, date_created, email_confirmed FROM users WHERE user_id = $1", id)

	u := User{}

	err := row.Scan(&u.ID, &email, &name, &u.Password, &u.EmailConfirmed, &u.Created, &u.EmailConfirmed)

	if email.Valid {
		u.Email = email.String
	} else {
		u.Email = ""
	}

	if name.Valid {
		u.Name = name.String
	} else {
		u.Name = ""
	}

	if err != nil {
		return User{}, err
	}

	// Convert time to local
	u.Created = u.Created.In(time.Local)

	return u, nil
}

func (us *UserService) ConfEmail(id string) error {
	_, err := us.db.Exec("UPDATE users SET email_confirmed = TRUE WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("Cannot change email_confirmed to true: %v", err)
	}
	return nil
}
