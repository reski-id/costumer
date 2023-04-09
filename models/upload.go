package models

import (
	"time"
)

type Upload struct {
	Model
	Name      string `json:"name" form:"name" binding:"required"`
	File      string `json:"file" form:"file" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
