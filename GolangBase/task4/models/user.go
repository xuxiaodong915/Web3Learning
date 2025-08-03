package models

import (
	"GolangBase/task4/utils/token"
	"errors"
	"html"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Email    string `gorm:"unique;not null"`
}


func LoginCheck(username, password string) (string, error) {
	var u User
	err := DB.Model(&User{}).Where("username=?", username).First(&u).Error
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword{
		return "", err
	}
	token ,err := token.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *User) SaveUser() (*User, error) {
	err := DB.Create(&u).Error
	if err != nil {
		log.Printf("saveUser err,%s", err)
		return &User{}, err
	}
	log.Printf("saveUser success,%+v", u)
	return u, nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("err,%s", err)
		return err
	}
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = string(hashedPassword)
	log.Printf("hashedPassword,%s", hashedPassword)
	return nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserByID(uid uint) (User,error) {
	var u User 
	if err := DB.First(&u,uid).Error;err != nil {
		return u,errors.New("user not found")
	}
	u.PrepareGive()
	return u, nil
}
