package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	// SignInTime    time.Time `gorm:"default: NULL"`
	// SignOutTime   time.Time `gorm:"default: NULL"`
	// HeartbeatTime time.Time `gorm:"default: NULL"`
	// Online        bool
	// ClientIp      string
	// ClientPort    string
}

// func FindUserById(id int) (*User, error) {
// 	var user User
// 	err := db.First(&user, id).Error
// 	return &user, err
// }

func FindUserByName(username string) (*User, error) {
	var user User
	// err := db.Where("username = ?", username).First(&user).Error
	err := db.First(&user, "username = ?", username).Error
	return &user, err
}

func FindUserByNameAndPwd(username string, password string) (*User, error) {
	var user User
	// err := db.Where("username = ?", username).First(&user).Error
	err := db.First(&user, "username = ? and password = ?", username, password).Error
	return &user, err
}

func CreateUser(user *User) error {
	return db.Create(user).Error
}

// func GetOnline(id int) {

// }
