package main

import (
	// "database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	TIMEOUT = 60
)

var (
	bot *tgbotapi.BotAPI
)

// U must create bot_tocken.go file, which include TOKEN variable in global package scope
func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(TOKEN)

	if err != nil {
		err.Error()
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = TIMEOUT

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		err.Error()
	}

	for update := range updates {
		msg := update.Message
		if msg == nil {
			continue
		}
		if msg.IsCommand() {
			command := msg.Command()

			switch command {
			case "regpi":
				regpi(msg)
			case "showpid":
				fmt.Println("nothing here")
			case "pidor":
				fmt.Println("nothing here")
			case "pidorstat":
				fmt.Println("nothing here")
			}
		}
	}

}
