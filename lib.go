package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
)

func createAvailableRecord(msg *tgbotapi.Message) {
	var available Available
	var user User

	userId := msg.From.ID
	groupId := msg.Chat.ID

	gdb.Where("groupId = ?", groupId).First(&available)
	gdb.Where("userId = ? AND groupId = ?", userId, groupId).First(&user)

	if available.Id == 0 {
		gdb.Create(&Available{
			GroupId: int(groupId),
			Flag:    true,
		})
	}
}

func createGroupRecord(msg *tgbotapi.Message) {
	var group Group
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)

	if group.Id == 0 {
		gdb.Create(&Group{
			GroupId: int(msg.Chat.ID),
			Title:   msg.Chat.Title,
			Name:    msg.Chat.UserName,
		})
	}
}

func checkIfUsernameChanged(msg *tgbotapi.Message) {
	var user User
	var group Group
	var reply tgbotapi.MessageConfig
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("userId = ? AND groupId = ?", msg.From.ID, group.Id).First(&user)
	castedUser := string(user.Username)
	if len(castedUser) > 0 && (castedUser[1:] != msg.From.UserName) {
		newUsername := msg.From.UserName
		user.Username = "@" + newUsername
		gdb.Model(&user).Update(User{Username: user.Username})
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Я помнил тебя под именем %s, запомню и новое имя %s",
			castedUser, "@"+newUsername))
		bot.Send(reply)
	}
}

func startSchedule() {
	gocron.Every(1).Day().At("12:00").Do(resetFlags)
	<-gocron.Start()
}

func resetFlags() {
	var user User
	var available Available
	gdb.Model(&available).Update("flag", true)
	gdb.Model(&user).Update(User{Quota: 6})
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn((max + 1)-min) + min
}

func checkIfPresenceUserNick(msg *tgbotapi.Message) {
	var user User

	userId := msg.From.ID
	groupId := msg.Chat.ID

	gdb.Where("userId = ? and groupId = ?", userId, groupId).First(&user)

	if user.Id > 0 && user.Usernick != fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName) {
		gdb.Model(&user).Where("userId = ?", userId).UpdateColumn("usernick",
			fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))
	}
}

func checkAdminAccess(msg *tgbotapi.Message, update tgbotapi.Update) (bool){
	var user User

	gdb.Where("userId = ? and groupId = ?", msg.From.ID, msg.Chat.ID).Find(&user)

	if user.IsAdmin == true {
		return true
	} else {
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Позарился на святое!"))
		reply.ReplyToMessageID = update.Message.MessageID
		bot.Send(reply)

		return false
	}
}
