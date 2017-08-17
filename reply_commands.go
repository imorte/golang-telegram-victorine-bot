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
			string(msg.Chat.ID), // это ваще то что надо, мне кажется что нет, но это не точно
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

func pidorStat(msg *tgbotapi.Message) {
	row, err := db.Query("SELECT pidor, score FROM pidors ORDER BY score DESC")

	if err != nil {
		err.Error()
	}

	var reply tgbotapi.MessageConfig
	var pidor string
	var score int
	var flag bool

	output := "Статистика:\n"
	for row.Next() {
		err = row.Scan(&pidor, &score)
		if err != nil {
			err.Error()
		}
		if score != 0 {
			flag = true
			output += fmt.Sprintf("%s: %d\n", pidor, score)
		}
	}

	if flag {
		reply = tgbotapi.NewMessage(msg.Chat.ID, output)
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Пидор дня еще ни разу не был выбран! /pidor")
	}

	bot.Send(reply)
}

func startQuiz(msg *tgbotapi.Message) {
	rows, err := db.Query("SELECT COUNT (*) FROM pidors")

	if err != nil {
		fmt.Printf("%s", err)
	}

	rowsCounted := checkCount(rows)

	println(rowsCounted)
	getWeather()

	//fmt.Println(cast(50, 1, 100, 1, 10))
}

func cast(x int, in_min int, in_max int, out_min int, out_max int) (int) {
	return (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err:= rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
