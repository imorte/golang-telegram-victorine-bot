package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"strconv"
	"time"
)

func regpi(msg *tgbotapi.Message, update tgbotapi.Update) {
	groupId := msg.Chat.ID

	var pidor Pidor
	gdb.Where("pidorId = ? AND wich_group = ?", msg.From.ID, groupId).First(&pidor)
	// testing us
	fmt.Println("--------- DEBUG --------")
	fmt.Println("-   ", pidor, "-")
	fmt.Println("--------- DEBUG  --------")
	// testing us END
	var reply tgbotapi.MessageConfig
	castedUser := string(pidor.Pidor)
	if pidor.ID == 0 {
		pidor.Pidor = "@" + msg.From.UserName
		pidor.PidorId = strconv.Itoa(int(msg.From.ID))
		pidor.WhichGroup = strconv.Itoa(int(msg.Chat.ID))
		pidor.Score = "0"
		gdb.Create(&pidor)
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprint("Ты зареган!"))
	} else if castedUser[1:] != msg.From.UserName {
		newUsername := msg.From.UserName
		pidor.Pidor = "@" + newUsername
		gdb.Model(&pidor).Update(Pidor{Pidor: pidor.Pidor})
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Я помнил тебя под именем %s, запомню и новое имя %s",
			castedUser, "@"+newUsername))
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprint("Эй, ты уже в игре!"))
	}

	reply.ReplyToMessageID = update.Message.MessageID
	bot.Send(reply)
}

func showpid(msg *tgbotapi.Message) {
	var pidors []Pidor
	gdb.Where("wich_group = ?", msg.Chat.ID).Find(&pidors)

	output := "Кандидаты в пидоры дня:\n"
	for _, i := range pidors {
		output += i.Pidor + "\n"
	}
	output += " Хочешь себя увидеть тут?\nЖми /regpi"
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, output))
}

func pidorStat(msg *tgbotapi.Message) {
	var pidors []Pidor
	var reply tgbotapi.MessageConfig
	gdb.Order("score desc").Find(&pidors)
	var flag bool

	output := "Статистика:\n"
	for _, i := range pidors {
		if i.Score != "0" {
			output += fmt.Sprintf("%s: %s\n", i.Pidor, i.Score)
			flag = true
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
	var pidors []Pidor
	gdb.Find(&pidors)

	rowsCounted := len(pidors)
	moscowWeather, oymyakonWeather := getWeather()
	averageWeather := (moscowWeather + oymyakonWeather) / 2

	calculatedWeather := cast(averageWeather, oymyakonWeather, moscowWeather, 1, rowsCounted)

	var thePidor int
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
