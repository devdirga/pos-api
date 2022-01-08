package model

import "time"

type (
	Student struct {
		Id        int       `gorm:"primary_key"`
		CreatedAt time.Time `json:"-"`
		UpdatedAt time.Time `json:"-"`
		Name      string    `gorm:"not null; type:varchar(100)"` // unique_index
	}
)
