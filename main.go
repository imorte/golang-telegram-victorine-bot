package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/jasonlvhit/gocron"
)

const (
	TIMEOUT       = 60
	DATABASE      = "sqlite3"
	DATABASE_NAME = "db.sqlite3"
)

var (
	bot *tgbotapi.BotAPI
	gdb *gorm.DB
)

// You must create bot_token.go file, which include TOKEN variable in global package scope
func init() {
	var err error

	gdb, err = gorm.Open(DATABASE, DATABASE_NAME)
	if err != nil {
		panic(err)
	}
	gdb.LogMode(true)

	gdb.AutoMigrate(
		&Users{},
		&Groups{},
		&Available{},
	)

	bot, err = tgbotapi.NewBotAPI(TOKEN)

	if err != nil {
		err.Error()
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

}

func main() {
	defer gdb.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = TIMEOUT

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		err.Error()
	}

	gocron.Every(1).Day().At("12:00").Do(resetAvailableFlag)


	for update := range updates {
		msg := update.Message
		if msg == nil {
			continue
		}
		if msg.IsCommand() {
			command := msg.Command()
			createGroupRecord(msg)
			createAvailableRecord(msg)
			checkIfUsernameChanged(msg)
			switch command {
				case "regpi":
					regpi(msg, update)
				case "showpid":
					showpid(msg)
				case "pidor":
					startQuiz(msg)
				case "pidorstat":
					pidorStat(msg)
			}
		}
	}

}
