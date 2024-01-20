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

	Name         string
	Email        sql.NullString
	PasswordHash sql.NullString

	Admin bool `gorm:"default:FALSE"`
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

type EmergencyContact struct {
	gorm.Model
	UserID uint
	User   User

	Name  string
	Phone string
}

type UserInfo struct {
	gorm.Model

	UserID uint
	User   User

	Address string

	Phone string
	Notes string

	ApprovedByID *uint
	ApprovedBy   *User
	ApprovedAt   *time.Time
}
