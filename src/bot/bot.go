package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"os"
)

type Bot struct {
	bot *telebot.Bot
}

func New() *Bot {
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
	return &Bot{tgBot}
}

func (b *Bot) SetWebhook(url string) {
	err := b.bot.SetWebhook(&telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: url,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic("can't set webhook")
	}
}

func (b *Bot) ProcessUpdate(u telebot.Update) {
	b.bot.ProcessUpdate(u)
}
