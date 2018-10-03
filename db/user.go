package db

import (
	"encoding/base64"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user that can login and make changes to the site
// configuration. Administrators can also manage users.
type User struct {
	ID       int64  `json:"id"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	IsAdmin  bool   `gorm:"not null" json:"is-admin"`
}

// GetName retrieves a descriptive name for groups of users.
func (u *User) GetName() string {
	return "users"
}

// GetID returns the unique identifier for the user.
func (u *User) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}

// SetID sets the unique identifier for the user.
func (u *User) SetID(id string) error {
	u.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}

// Authenticate hashes the password and compares it to the value stored in the
// database. An error is returned if the values do not match.
func (u *User) Authenticate(password string) error {
	h, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(h, []byte(password))
}

// SetPassword salts and hashes the user's password. It does not store the new
// value in the database.
func (u *User) SetPassword(password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	u.Password = base64.StdEncoding.EncodeToString(h)
	return nil
}
