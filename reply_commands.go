package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"time"
)

func regpi(msg *tgbotapi.Message, update tgbotapi.Update) {
	groupId := msg.Chat.ID
	var (
		id sql.NullInt64
		user sql.NullString
		group sql.NullInt64
		score sql.NullInt64
		userId sql.NullInt64
	)

	err := db.QueryRow("SELECT * FROM pidors WHERE pidorId = ? AND wich_group = ? ",
		msg.From.ID, groupId).Scan(&id, &user, &group, &score, &userId)

	if err != nil {
		err.Error()
	}


	var reply tgbotapi.MessageConfig
	castedUser := string(user.String)
	if !id.Valid {
		_, err = db.Exec(
			"INSERT INTO pidors (pidor, pidorId, wich_group, score) VALUES (?, ?, ?, ?)",
			"@" + msg.From.UserName,
			msg.From.ID,
			msg.Chat.ID,
			0,
		)
		if err != nil {
			err.Error()
		}
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ты регнулся, @%s", msg.From.UserName))
	} else if castedUser[1:] != msg.From.UserName {
		newUsername := msg.From.UserName
		_, err := db.Exec("UPDATE pidors SET pidor = ? where pidorId = ? and wich_group = ?", "@" + newUsername, userId, msg.Chat.ID)
		if err != nil {
			err.Error()
		}
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Я помнил тебя под именем %s, запомню и новое имя %s",
			castedUser, "@" + newUsername))
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprint("Эй, ты уже в игре!"))
	}

	reply.ReplyToMessageID = update.Message.MessageID
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
	output += " Хочешь себя увидеть тут?\nЖми /regpi"
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
	var thePidor int

	rows, err := db.Query("SELECT COUNT (*) FROM pidors")
	if err != nil {
		fmt.Printf("%s", err)
	}

	rowsCounted := checkCount(rows)

	moscowWeather, oymyakonWeather := getWeather()
	averageWeather := (moscowWeather + oymyakonWeather) / 2

	calculatedWeather := cast(averageWeather, oymyakonWeather, moscowWeather, 1, rowsCounted)

	if calculatedWeather > rowsCounted/2 {
		thePidor = random(1, calculatedWeather/2)
	} else {
		thePidor = random(calculatedWeather, rowsCounted)
	}

	println(thePidor)

}

func cast(x int, inMin int, inMax int, outMin int, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
