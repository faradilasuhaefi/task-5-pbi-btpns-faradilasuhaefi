package models

import (
	"time"

	"final-project-pbi-btpns/helpers"
)

type User struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey; not null; autoIncrement"`
	Username  string    `json:"username" gorm:"type:varchar(30)"`
	Email     string    `json:"email" gorm:"unique; type:varchar(30)"`
	Password  string    `json:"password"`
	Photos    []Photo   `gorm:"foreignKey:UserID; constarint:OnUpdate:CASCADE, onDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (user *User) HashPassword(password string) error {
	hash, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hash
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	result, err := helpers.ComparePassword(providedPassword, user.Password)
	if !result {
		return err
	}
	return nil
}
