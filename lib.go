package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
	"time"
	"math/rand"
)

func createAvailableRecord(msg *tgbotapi.Message) {
	var available Available
	var group Groups
	var user Users
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("groupId = ?", group.Id).First(&available)
	gdb.Where("userId = ? AND groupId = ?", msg.From.ID, group.Id).First(&user)

	if available.Id == 0 {
		gdb.Create(&Available{
			GroupId: group.Id,
			Flag: true,
		})
	}
}

func createGroupRecord(msg *tgbotapi.Message) {
	var group Groups
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)

	if group.Id == 0 {
		gdb.Create(&Groups{
			GroupId: int(msg.Chat.ID),
			Title: msg.Chat.Title,
			Name: msg.Chat.UserName,
		})
	}
}

func checkIfUsernameChanged(msg *tgbotapi.Message) {
	var user Users
	var group Groups
	var reply tgbotapi.MessageConfig
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("userId = ? AND groupId = ?", msg.From.ID, group.Id).First(&user)
	castedUser := string(user.Username)
	if len(castedUser) > 0 && (castedUser[1:] != msg.From.UserName) {
		newUsername := msg.From.UserName
		user.Username = "@" + newUsername
		gdb.Model(&user).Update(Users{Username: user.Username})
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Я помнил тебя под именем %s, запомню и новое имя %s",
			castedUser, "@"+newUsername))
		bot.Send(reply)
	}
}

func cast(x int, inMin int, inMax int, outMin int, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
