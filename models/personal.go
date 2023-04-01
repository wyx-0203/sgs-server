package models

import (
	"math/rand"
	"strconv"
	"time"
)

// import "gorm.io/gorm"

type Personal struct {
	UserID    uint `gorm:"primarykey"`
	Name      string
	Character string
	Win       uint
	Lose      uint
}

func FindPersonal(id uint) (*Personal, error) {
	var personal Personal
	err := db.First(&personal, id).Error
	return &personal, err
}

func CreatePersonal(id uint) error {
	// 生成初始昵称
	rand.Seed(time.Now().UnixNano())
	name := "用户" + strconv.Itoa(rand.Intn(999999))

	personal := &Personal{
		UserID:    id,
		Name:      name,
		Character: "100101",
	}

	return db.Create(personal).Error
}

func Rename(p *Personal, name string) error {
	p.Name = name
	return db.Save(p).Error
}

func ChangeCharacter(p *Personal, character string) error {
	p.Character = character
	return db.Save(p).Error
}
