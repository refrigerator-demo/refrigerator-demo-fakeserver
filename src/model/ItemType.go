package model

import (
	"time"
)

type ItemType struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	TypeEnum  int       `gorm:"not null;unique" json:"type_enum"`
	TypeTitle string    `gorm:"size:255;not null;unique" json:"type_title"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
