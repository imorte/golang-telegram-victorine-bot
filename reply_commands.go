package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"time"
)

func regpi(msg *tgbotapi.Message, update tgbotapi.Update) {
	var reply tgbotapi.MessageConfig
	var user Users
	var group Groups
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("userId = ? AND groupId = ?", msg.From.ID, group.Id).First(&user)


	if len(msg.From.UserName) == 0 {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Сначала добавь ник, а потом играй!")
	} else if user.Id == 0 {
		gdb.Where("groupId = ?", msg.Chat.ID).First(&group)

		user.Username = "@" + msg.From.UserName
		user.UserId = msg.From.ID
		user.GroupId = group.Id
		user.Score = 0
		gdb.Create(&user)
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ты регнулся, %s", user.Username))
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprint("Эй, ты уже в игре!"))
	}

	reply.ReplyToMessageID = update.Message.MessageID
	bot.Send(reply)
}

func showpid(msg *tgbotapi.Message) {
	var group Groups
	var users []Users
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("groupId = ?", group.Id).Find(&users)

	if len(users) != 0 {
		output := "Кандидаты в пидоры дня:\n"
		for _, i := range users {
			output += i.Username + "\n"
		}
		output += "Хочешь себя увидеть тут?\nЖми /regpi"
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, output))
	} else {
		output := "Пидоров нет! Будь первым! Жми /regpi"
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, output))
	}
}

func pidorStat(msg *tgbotapi.Message) {
	var group Groups
	var users []Users
	var reply tgbotapi.MessageConfig
	var flag bool
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("groupId = ?", group.Id).Order("score desc").Find(&users)

	output := "Статистика:\n"
	for _, i := range users {
		if i.Score != 0 {
			output += fmt.Sprintf("%s: %d\n", i.Username, i.Score)
			flag = true
		}
	}

	if flag {
		reply = tgbotapi.NewMessage(msg.Chat.ID, output)
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Пидор дня еще ни разу не был выбран! Жми /pidor")
	}

	bot.Send(reply)
}

func startQuiz(msg *tgbotapi.Message) {
	firstPhrases := []string {
		"Инициализирую поиск пидора дня...",
		"Внимание, ищу пидора!",
		"Ну-ка дай-ка...",
		"Такс, кто тут у нас мало каши ел?",
		"Инициализация.Поиск.",
	}

	secondPhrases := []string {
		"Кажется я что-то вижу!",
		"Не может быть!",
		"Пожалуй препроверю...",
		"Найден!",
		"Прям по Бабичу!",
		"Как предсказал Великий Мейстер...",
	}

	var reply tgbotapi.MessageConfig
	var theUser int
	var users []Users
	var group Groups
	var winner Users
	var winnerScore int
	var available Available
	gdb.Where("groupId = ?", msg.Chat.ID).First(&group)
	gdb.Where("groupId = ?", group.Id).Find(&users)
	gdb.Where("groupId = ?", group.Id).First(&available)

	rowsCounted := len(users)
	if rowsCounted == 0 {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Нет участников! Жми /regpi")
		bot.Send(reply)
	} else {
		if available.Flag {
			lenOfCurrentUsers := len(users)
			theUser = random(0, lenOfCurrentUsers - 1)


			println()
			println(theUser)
			println()


			reply = tgbotapi.NewMessage(msg.Chat.ID, firstPhrases[random(0, len(secondPhrases) - 1)])
			bot.Send(reply)
			time.Sleep(time.Second * 2)
			reply = tgbotapi.NewMessage(msg.Chat.ID, secondPhrases[random(0, len(firstPhrases) - 1)])
			bot.Send(reply)
			gdb.Where("id = ? and groupId = ?", theUser, group.Id).First(&winner)
			winnerScore = winner.Score + 1
			gdb.Model(&users).Where("id = ?", winner.Id).UpdateColumn("score", winnerScore)
			time.Sleep(time.Second * 2)
			reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ага! 🎉🎉🎉 Сегодня пидор - %s", winner.Username))
			bot.Send(reply)
			gdb.Model(&available).Where("groupId = ?", group.Id).Update("flag", false)
			gdb.Model(&available).Where("groupId = ?", group.Id).Update("userId", winner.Id)
		} else {
			var currentUser Users
			gdb.Where("id = ?", available.UserId).First(&currentUser)
			reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("🎉Сегодня у нас уже есть победитель - %s🎉", currentUser.Username))
			bot.Send(reply)
		}
	}
}

