package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"time"
	//"os/user"
)

func regpi(msg *tgbotapi.Message, update tgbotapi.Update) {
	var reply tgbotapi.MessageConfig
	var user User

	userId := msg.From.ID
	groupId := msg.Chat.ID

	gdb.Where("userId = ? AND groupId = ?", userId, groupId).First(&user)


	if len(msg.From.UserName) == 0 {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "A girl has no name.")
	} else if user.Id == 0 {
		user.Username = "@" + msg.From.UserName
		user.UserId = msg.From.ID
		user.GroupID = int(groupId)
		user.Score = 0
		user.Usernick = fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName)
		user.Quota = 6
		gdb.Create(&user)
		gdb.Model(&user).Update(User{Quota: 6})
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ты регнулся, [%s](tg://user?id=%d)\n", user.Usernick, user.UserId))
		reply.ParseMode = tgbotapi.ModeMarkdown
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprint("Эй, ты уже в игре!"))
	}

	reply.ReplyToMessageID = update.Message.MessageID
	bot.Send(reply)
}

func showpid(msg *tgbotapi.Message) {
	var users []User

	groupId := msg.Chat.ID

	gdb.Where("groupId = ?", groupId).Find(&users)

	if len(users) != 0 {
		output := "Кандидаты в пидоры дня:\n"
		for _, i := range users {
			if len(i.Usernick) > 0 {
				output += fmt.Sprintf("%s\n", i.Usernick)
			} else {
				output += fmt.Sprintf("%s\n", i.Username[1:])
			}
		}
		output += "Хочешь себя увидеть тут?\nЖми /regpi"
		reply := tgbotapi.NewMessage(msg.Chat.ID, output)
		bot.Send(reply)
	} else {
		output := "Пидоров нет! Будь первым! Жми /regpi"
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, output))
	}
}

func pidorStat(msg *tgbotapi.Message) {
	titles := map[int]string{
		1 : "Пидоратор",
		2 : "Пидороль",
		3 : "Герцопидор",
		4 : "Пиркиз",
		5 : "Пидорон",
	}

	var users []User
	var reply tgbotapi.MessageConfig
	var flag bool
	var currentUserName string

	groupId := msg.Chat.ID
	counter := 0
	var titlesCounter int

	gdb.Where("groupId = ?", groupId).Order("score desc").Find(&users)

	output := "Статистика (первые 5):\n"
	titlesCounter = 1
	for _, i := range users {
		if i.Score != 0 {
			if len(i.Usernick) > 0 {
				currentUserName = i.Usernick
			} else {
				currentUserName = i.Username[1:]
			}

			output += fmt.Sprintf("[%s](tg://user?id=%d) - %d (%s)\n", currentUserName, i.UserId, i.Score, titles[titlesCounter])
			titlesCounter++
			flag = true

			if counter == 4 {
				break
			} else {
				counter++
			}
		}
	}

	if flag {
		reply = tgbotapi.NewMessage(msg.Chat.ID, output)
		reply.ParseMode = tgbotapi.ModeMarkdown
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
	var theUser User
	var users []User
	var randomUser int
	var currentUserName string
	var winner User
	var winnerScore int
	var available Available

	groupId := msg.Chat.ID

	gdb.Where("groupId = ?", groupId).Find(&users)
	gdb.Where("groupId = ?", groupId).First(&available)

	rowsCounted := len(users)
	if rowsCounted == 0 {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Нет участников! Жми /regpi")
		bot.Send(reply)
	} else {
		if available.Flag {
			lenOfCurrentUsers := len(users)
			if lenOfCurrentUsers == 1 {
				randomUser = 0
			} else {
				randomUser = random(0, lenOfCurrentUsers - 1)
			}

			gdb.Where("id = ?", users[randomUser].Id).First(&winner)

			reply = tgbotapi.NewMessage(msg.Chat.ID, firstPhrases[random(0, len(firstPhrases) - 1)])
			bot.Send(reply)
			time.Sleep(time.Second * 2)
			reply = tgbotapi.NewMessage(msg.Chat.ID, secondPhrases[random(0, len(secondPhrases) - 1)])
			bot.Send(reply)
			gdb.Where("id = ? and groupId = ?", theUser, groupId).First(&winner)
			winnerScore = winner.Score + 1
			gdb.Model(&users).Where("id = ?", winner.Id).UpdateColumn("score", winnerScore)
			time.Sleep(time.Second * 2)
			if len(winner.Usernick) > 0 {
				currentUserName = winner.Usernick
			} else {
				currentUserName = winner.Username
			}
			reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ага! 🎉🎉🎉 Сегодня пидор - [%s](tg://user?id=%d)", currentUserName, winner.UserId))
			reply.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(reply)
			gdb.Model(&available).Where("groupId = ?", groupId).Update("flag", false)
			gdb.Model(&available).Where("groupId = ?", groupId).Update("userId", winner.Id)
		} else {
			var currentUser User
			gdb.Where("id = ?", available.UserId).First(&currentUser)
			if len(currentUser.Usernick) > 0 {
				currentUserName = currentUser.Usernick
			} else {
				currentUserName = currentUser.Username[1:]
			}
			reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("🎉Сегодня у нас уже есть победитель - [%s](tg://user?id=%d)🎉", currentUserName, currentUser.UserId))
			reply.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(reply)
		}
	}
}

func kekogen(msg *tgbotapi.Message) {
	var reply tgbotapi.MessageConfig
	var user User

	userId := msg.From.ID
	groupId := msg.Chat.ID

	gdb.Where("userId = ? and groupId = ?", userId, groupId).First(&user)
	currentQuota := user.Quota

	if user.Id > 0 {
		if currentQuota > 1 {
			vowels := []string {
				"а", "о", "и", "е", "у",  "я",
			}
			consonants := []string {
				"в", "д", "к", "л", "м", "н", "п", "р", "с", "т", "ф", "х", "ш",
			}
			result := "кек"

			for x:= 0; x < 5; x++ {
				if x % 2 == 0 {
					result += vowels[random(0, len(vowels) - 1)]
				} else {
					result += consonants[random(0, len(consonants) - 1)]
				}
			}

			gdb.Model(&user).Update(User{Quota: currentQuota - 1})
			reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf(result))
		} else {
			reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Твои кеки на сегодня кончились!"))
		}
	} else {
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Жми /regpi чтобы зарегистрироваться и получить свои кеки!"))
	}

	bot.Send(reply)
}