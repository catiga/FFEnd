package model

type MysticalArt struct {
	BaseModel
	Name  string
	Brief string
}

func (MysticalArt) TableName() string {
	return "spw_mysart"
}

type Catalog struct {
	BaseModel
	Name   string
	Seq    int
	Parent uint64
}

func (Catalog) TableName() string {
	return "spw_catag"
}

type CharacterPos struct {
	Id       uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	CharId   uint64
	TypeLan  string
	TypeCode string
	TypeCat  string
	Flag     int
}

func (CharacterPos) TableName() string {
	return "spw_char_position"
}
