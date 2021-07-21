package model

import (
	"time"

	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Inventory struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	/*InventoryType InventoryType `gorm:"foreignKey:type_id;" json:"inventory_type"`*/
	UserId      uint64    ``
	Title       string    `gorm:"size:255;not null;" json:"title"`
	Description string    `gorm:"type:text;not null;" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (inventory *Inventory) Initialize() {
	inventory.ID = 0
	inventory.UserId = 0
	inventory.Title = html.EscapeString(strings.TrimSpace(inventory.Title))
	inventory.Description = html.EscapeString(strings.TrimSpace(inventory.Description))
	inventory.CreatedAt = time.Now()
	inventory.UpdatedAt = time.Now()
}

func (inventory *Inventory) DoBeforeSave(userId uint64) {
	inventory.UserId = userId
}

func (inventory *Inventory) FindInventoryById(db* gorm.DB, userId uint64) (*[]Inventory, error) {
	
	inventories := []Inventory{}

	err := db.Debug().Model(&Inventory{}).Where("user_id = ?", userId).Limit(100).Find(&inventories).Error
	if nil != err {
		return &[]Inventory{}, err
	}

	return &inventories, err
}

func (inventory *Inventory) CreateInventory(db* gorm.DB) error {
	// Begin Transaction
	tx := db.Debug().Begin()

	// recover defer
	defer func() {
		if r := recover(); nil != r {
			tx.Rollback()
		}
	}()

	// DB Logic Start
	err := tx.Create(&inventory).Error
	if nil != err {
		tx.Rollback()
		return err
	}

	// DB Logic End

	// Commit Transaction
	err = tx.Commit().Error
	if nil != err {
		tx.Rollback()

		return err
	}
	// Return Value
	return nil
}