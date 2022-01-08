package handler

import "github.com/jinzhu/gorm"

type (
	Handler struct {
		GormDB *gorm.DB
	}
)

const (
	Key = "secret"
)
