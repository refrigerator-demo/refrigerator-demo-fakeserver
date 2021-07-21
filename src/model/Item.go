package model

import "time"

type Item struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	/*Inventory   Inventory `gorm:"foreignKey:InventoryRefer;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;;" json:"storage"`*/
	InventoryId uint64    ``
	TypeId      uint64    ``
	Category    uint32    `gorm:"not null;" json:"category"`
	Count       uint32    `gorm:"not null;" json:"count"`
	Title       string    `gorm:"size:255;not null;unique" json:"title"`
	Description string    `gorm:"type:text;not null;" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
