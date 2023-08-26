package model

import "time"

type BaseModel struct {
	Id   uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Code string
	Lan  string
	Flag int
}

type Character struct {
	BaseModel
	CharName       string
	CharAvatar     string
	CharInfo       string
	CharBirth      string
	CharAge        string
	CharGender     string
	CharPlace      string
	CharFullBody   string
	CharProfile    string
	CharNatureCode string
	CharRegion     string
}

func (Character) TableName() string {
	return "spw_character"
}

type Method struct {
	BaseModel
	Name string
	Flag int
}

func (Method) TableName() string {
	return "spw_method"
}

type CharBack struct {
	BaseModel
	CharId  uint64
	Role    string
	Prompt  string
	Seq     int
	AddTime *time.Time
}

func (CharBack) TableName() string {
	return "spw_char_background"
}
