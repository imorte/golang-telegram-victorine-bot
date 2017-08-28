package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
)

func checkForSchedulePresence(msg *tgbotapi.Message) {
	var ava Available
	// here 1 result from query????
	gdb.Where("group_telega = ?", msg.Chat.ID).First(&ava)
	// err := db.QueryRow("SELECT id from available WHERE group_telega = ?", msg.Chat.ID).Scan(&available)

	if ava.Id == 0 {
		gdb.Create(&Available{
			GroupTelega: strconv.Itoa(int(msg.Chat.ID)),
			Flag:        "0",
			Current:     "",
		})
	}
}
