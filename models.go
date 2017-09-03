package main

type Users struct {
	Id       int    `gorm:"primary_key;column:id"`
	Username string `gorm:"type:TEXT;column:username"`
	GroupId  int      `gorm:"column:groupId"`
	Score    int      `gorm:"column:score"`
	UserId   int      `gorm:"column:userId"`
}

func (Users) TableName() string {
	return "users"
}

type Groups struct {
	Id      int `gorm:"primary_key;column:id"`
	GroupId int `gorm:"column:groupId"`
	Title   string `gorm:"type:TEXT;column:title"`
	Name    string `gorm:"type:TEXT;column:name"`
}

func (Groups) TableName() string {
	return "groups"
}

type Available struct {
	Id      int    `gorm:"primary_key;column:id"`
	GroupId int    `gorm:"column:groupId"`
	Flag    bool   `gorm:"column:flag"`
	UserId  int    `gorm:"column:userId"`
}

func (Available) TableName() string {
	return "available"
}
