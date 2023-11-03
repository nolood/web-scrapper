package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"web-scrapper/pkg/environment"
)

var bot *tgbotapi.BotAPI

func InitTelegramBot() {
	var err error

	token := environment.GetValue("TG_TOKEN")

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	//bot.Debug = true

}

func Send(msg string) {
	users := []int64{665209714, 2117838522}

	for _, user := range users {
		message := tgbotapi.NewMessage(user, msg)
		bot.Send(message)
	}
}
