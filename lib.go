package main

import (
	"gopkg.in/telegram-bot-api.v4"
)

func checkForSchedulePresence(msg *tgbotapi.Message) {
	var available int
	err := db.QueryRow("SELECT id from available WHERE group_telega = ?", msg.Chat.ID).Scan(&available)

	if err != nil {
		err.Error()
	}

	if available == 0 {
		_, err = db.Exec("INSERT INTO available (group_telega, flag, current) VALUES (?, ?, ?)",
			msg.Chat.ID, 0, "")
	}
}