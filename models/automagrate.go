package models

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, secret string) {
	// Register all the models here
	err := db.AutoMigrate(
		&User{},
	)
	if err != nil {
		logrus.Error(err)
	}

	// Bootstrap
	user := &User{
		Admin: true,
	}
	user_attr := &User{
		Email: S("admin@admin.com"),
	}
	user_attr.SetPassword("admin", secret)
	db.Where(user).
		Attrs(user_attr).FirstOrCreate(user)
}

func S(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
