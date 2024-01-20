package models

import (
	"database/sql"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uint

	Name         string
	Email        sql.NullString
	PasswordHash sql.NullString

	Admin bool `gorm:"default:FALSE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) SetPassword(password, secret string) error {
	passwordHash, err := scrypt.Key(
		[]byte(password),
		[]byte(secret),
		32768, 8, 1, 32,
	)
	u.PasswordHash = S(base64.StdEncoding.EncodeToString(passwordHash))
	return err
}
