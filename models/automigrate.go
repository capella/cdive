package models

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AutoMigrate(db *gorm.DB, secret string) {
	// Register all the models here
	err := db.AutoMigrate(
		&User{},
		&EmergencyContact{},
		&UserInfo{},
		&Event{},
	)
	if err != nil {
		logrus.Error(err)
	}

	// Bootstrap
	user := &User{
		Admin: true,
	}
	user_attr := &User{
		Email: "admin@admin.com",
	}
	user_attr.SetPassword("admin", secret)
	db.Where(user).
		Attrs(user_attr).FirstOrCreate(user)

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&EventsCategory{Name: "Diving Program"})
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&EventsCategory{Name: "Training"})
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&EventsCategory{Name: "Social"})
}

func S(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
