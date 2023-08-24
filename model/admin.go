package model

import (
	"time"
)

type AdminUser struct {
	Id       uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	UserName string
	Password string
}

func (AdminUser) TableName() string {
	return "admin_user"
}

type AdminToken struct {
	Id         uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	AdminId    uint64
	Token      string
	ExpireTime *time.Time
	CreateTime *time.Time
	Flag       int
}

func (AdminToken) TableName() string {
	return "admin_token"
}
