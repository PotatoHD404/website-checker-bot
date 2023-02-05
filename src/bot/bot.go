package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"os"
)

func New() *telebot.Bot {
	settings := telebot.Settings{
		Token:       os.Getenv("BOT_TOKEN"),
		Synchronous: true,
		Verbose:     true,
	}
	newTgBot, err := telebot.NewBot(settings)
	if err != nil {
		fmt.Println(err)
		panic("can't create bot")
	}
	tgBot := newTgBot
	tgBot.Handle(telebot.OnText, func(m *telebot.Message) {
		message := m.Text
		_, err := tgBot.Send(m.Sender, message)
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})
	return tgBot
}
