package model

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func MakeHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) DoBeforeSave() error {
	hashedPassword, err := MakeHash(user.Password)
	if nil != err {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *User) Initialize() {
	user.ID = 0
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate() error {
	return nil
}

func (user *User) SaveUser(db *gorm.DB) (*User, error) {

	// Begin Transaction
	tx := db.Debug().Begin()

	// recover defer
	defer func() {
		if r := recover(); nil != r {
			tx.Rollback()
		}
	}()

	// DB Logic Start
	err := tx.Create(&user).Error
	if nil != err {
		return &User{}, err
	}

	// DB Logic End

	// Commit Transaction
	err = tx.Commit().Error

	// Return Value
	return user, err
}

func (user *User) FindUserByID(db *gorm.DB, uid uint64) (*User, error) {

	err := db.Debug().Model(User{}).Where("id = ?", uid).Take(&user).Error

	if nil != err {
		return &User{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return user, err
}

func (user *User) UpdateUser(db *gorm.DB, uid uint64) (*User, error) {
	err := user.DoBeforeSave()
	if nil != err {
		log.Fatal(err)
	}

	// Begin Transaction
	tx := db.Debug().Begin()

	// recover defer
	defer func() {
		if r := recover(); nil != r {
			tx.Rollback()
		}
	}()

	// DB Logic Start
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  user.Password,
			"name":      user.Name,
			"email":     user.Email,
			"update_at": time.Now(),
		},
	)

	if nil != db.Error {
		return &User{}, db.Error
	}

	// DB Logic End

	// Commit Transaction
	err = tx.Commit().Error

	// Return Value
	return user, err
}
