package main

type Pidor struct {
	ID         int    `gorm:"primary_key;column:id"`
	Pidor      string `gorm:"type:TEXT;column:pidor"`
	WhichGroup string `gorm:"type:TEXT;column:wich_group"`
	Score      string `gorm:"type:TEXT;column:score"`
	PidorId    string `gorm:"type:TEXT;column:pidorId"`
}

func (Pidor) TableName() string {
	return "pidors"
}

type Group struct {
	Id   int `gorm:"column:id"`
	Name int `gorm:"column:name"`
}

func (Group) TableName() string {
	return "groups"
}

type Available struct {
	Id          int    `gorm:"primary_key;column:id"`
	GroupTelega string `gorm:"column:group_telega"`
	Flag        string `gorm:"column:flag"`
	Current     string `gorm:"column:current"`
}

func (Available) TableName() string {
	return "available"
}
