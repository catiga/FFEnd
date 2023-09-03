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
