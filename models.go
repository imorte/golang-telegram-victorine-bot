package main

type User struct {
	Id       int    `gorm:"primary_key;column:id"`
	Username string `gorm:"type:TEXT;column:username"`
	Usernick string `gorm:"type:TEXT;column:usernick"`
	GroupID  int    `gorm:"column:groupId"`
	Score    int    `gorm:"column:score"`
	UserId   int    `gorm:"column:userId"`
	Quota    int    `gorm:"column:quota"`
	IsAdmin  bool   `gorm:"column:isAdmin"`
	DisableNotify  bool   `gorm:"column:disable_notify"`
}

func (User) TableName() string {
	return "users"
}

type Group struct {
	Id      int    `gorm:"primary_key;column:id"`
	GroupId int    `gorm:"column:groupId"`
	Title   string `gorm:"type:TEXT;column:title"`
	Name    string `gorm:"type:TEXT;column:name"`
}

func (Group) TableName() string {
	return "groups"
}

type Available struct {
	Id      int  `gorm:"primary_key;column:id"`
	GroupId int  `gorm:"column:groupId"`
	Flag    bool `gorm:"column:flag"`
	UserId  int  `gorm:"column:userId"`
}

func (Available) TableName() string {
	return "available"
}
