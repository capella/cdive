package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EventsCategory struct {
	gorm.Model

	Name      string
	ShortName string `gorm:"unique"`

	Description string
}

func (e *EventsCategory) BeforeCreate(tx *gorm.DB) (err error) {
	e.ShortName = strings.ReplaceAll(strings.ToLower(e.Name), " ", "-")
	return
}

type Event struct {
	gorm.Model

	Name        string
	Description *string
	Location    *string

	CategoryID uint
	Category   EventsCategory

	Start time.Time
	End   *time.Time

	AvalibleSpaces *uint
	Owners         []User `gorm:"many2many;"`
	Participants   []User `gorm:"many2many;"`
	WaitingList    []User `gorm:"many2many;"`

	MinimumQualifications *string
	DepositCode           *string
	Cost                  *float64
}

func (e *Event) BeforeSave(tx *gorm.DB) (err error) {
	if e.Name == "" {
		err = errors.New("can't save invalid data")
	}

	return
}
