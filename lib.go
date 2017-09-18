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
	var group Groups
	var user Users
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("groupId = ?", group.Id).First(&available)
	gdb.Where("userId = ? AND groupId = ?", msg.From.ID, group.Id).First(&user)

	if available.Id == 0 {
		gdb.Create(&Available{
			GroupId: group.Id,
			Flag:    true,
		})
	}
}

func createGroupRecord(msg *tgbotapi.Message) {
	var group Groups
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)

	if group.Id == 0 {
		gdb.Create(&Groups{
			GroupId: int(msg.Chat.ID),
			Title:   msg.Chat.Title,
			Name:    msg.Chat.UserName,
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

func startSchedule() {
	gocron.Every(1).Day().At("12:00").Do(resetFlags)
	<-gocron.Start()
}

func resetFlags() {
	var user Users
	var available Available
	gdb.Model(&available).Update("flag", true)
	gdb.Model(&user).Update(Users{Quota: 6})
}

func cast(x int, inMin int, inMax int, outMin int, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn((max + 1)-min) + min
}

//func digitToWord(digit string) string {
//	var result string
//	_ := map[int]string {
//		1: "один",
//		2: "два",
//		3: "три",
//		4: "четыре",
//		5: "пять",
//		6: "шесть",
//		7: "семь",
//		8: "восемь",
//		9: "девять",
//	}
//
//
//	return result
//}