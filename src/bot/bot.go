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

	tgBot.Handle(telebot.OnPhoto, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle photos yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnDocument, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle documents yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnSticker, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle stickers yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnAudio, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle audio yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnVoice, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle voice yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnVideo, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle video yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnVideoNote, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle video notes yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnContact, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle contacts yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnLocation, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle locations yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnVenue, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle venues yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnPoll, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle polls yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnDice, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle dice yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	tgBot.Handle(telebot.OnGame, func(m *telebot.Message) {
		_, err := tgBot.Send(m.Sender, "I can't handle games yet")
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

func (b *Bot) Send(to telebot.Recipient, what interface{}, options ...interface{}) (*telebot.Message, error) {
	return b.bot.Send(to, what, options...)
}
