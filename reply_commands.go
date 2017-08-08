package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
)

func regpi(msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ты зарегистрирован %s", msg.From.UserName))
	reply.ReplyToMessageID = msg.MessageID
	bot.Send(reply)
}
