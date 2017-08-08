package main

import (
	"database/sql"
	"fmt"

	"gopkg.in/telegram-bot-api.v4"

	_ "github.com/mattn/go-sqlite3"
)

func regpi(msg *tgbotapi.Message) {
	var result sql.NullInt64
	row := db.QueryRow(
		"SELECT id FROM pidors WHERE pidor=?",
		"@"+msg.From.UserName,
	)
	err := row.Scan(&result)
	if err != nil {
		err.Error()
	}

	var reply tgbotapi.MessageConfig

	if !result.Valid {
		_, err = db.Exec(
			"INSERT INTO pidors (pidor, wich_group, score) VALUES (?, ?, ?)",
			"@"+msg.Chat.UserName,
			string(msg.Chat.ID),
			0,
		)
		if err != nil {
			err.Error()
		}
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ты регнулся %s", msg.From.UserName))
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Ты уже зарегистрирован")
	}

	bot.Send(reply)
}

func showpid(msg *tgbotapi.Message) {
	row, err := db.Query("SELECT pidor FROM pidors")
	if err != nil {
		err.Error()
	}

	output := "Кандидаты в пидоры дня:\n"
	var pidorName string
	for row.Next() {
		err = row.Scan(&pidorName)

		if err != nil {
			err.Error()
		}
		output += pidorName + "\n"
	}
	output += " Хочешь себя увидеть тут? \nЖми /regpi"
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, output))

}
