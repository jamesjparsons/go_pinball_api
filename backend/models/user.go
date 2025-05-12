package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `swaggerignore:"true"`
	Email      string   `json:"email" gorm:"unique"`
	Password   string   `json:"-"` // The "-" tag means this field won't be included in JSON
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Leagues    []League `gorm:"foreignKey:OwnerID"`
}

// SetPassword hashes the password and stores it
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
