package bot

import (
	"fmt"
	"os"
	. "website-checker-bot/bot/commands"
	. "website-checker-bot/bot/commands/admin"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/commands/subscription"
	. "website-checker-bot/bot/commands/website"

	"gopkg.in/tucnak/telebot.v2"
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
	tgBot, err := telebot.NewBot(settings)
	if err != nil {
		fmt.Println(err)
		panic("can't create bot")
	}

	return &Bot{tgBot}
}

func (b *Bot) GetBot() *telebot.Bot {
	return b.bot
}

func (b *Bot) Init(e *Env) {
	b.bot.Handle(telebot.OnText, func(m *telebot.Message) {
		message := m.Text
		switch message {
		case "/start":
			HandleStart(m, e)
			break
		case "/help":
			HandleHelp(m, e)
			break
		case "/add_website":
			HandleAddWebsite(m, e)
			break
		case "/get_websites":
			HandleGetWebsites(m, e)
			break
		case "/delete_website":
			HandleDeleteWebsite(m, e)
			break
		case "/add_admin":
			HandleAddAdmin(m, e)
			break
		case "/get_admins":
			HandleGetAdmins(m, e)
			break
		case "/delete_admin":
			HandleDeleteAdmin(m, e)
			break
		case "/subscribe":
			HandleSubscribe(m, e)
			break
		case "/unsubscribe":
			HandleUnsubscribe(m, e)
			break
		case "/get_subscriptions":
			HandleGetSubscriptions(m, e)
			break

		default:
			_, err := b.bot.Send(m.Sender, "I don't understand you")
			if err != nil {
				fmt.Println(err)
				panic("can't send message")
			}
			break
		}
	})

	b.bot.Handle(telebot.OnPhoto, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle photos yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnDocument, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle documents yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnSticker, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle stickers yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnAudio, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle audio yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnVoice, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle voice yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnVideo, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle video yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnVideoNote, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle video notes yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnContact, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle contacts yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnLocation, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle locations yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnVenue, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle venues yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnPoll, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle polls yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnDice, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle dice yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})

	b.bot.Handle(telebot.OnGame, func(m *telebot.Message) {
		_, err := b.bot.Send(m.Sender, "I can't handle games yet")
		if err != nil {
			fmt.Println(err)
			panic("can't send message")
		}
	})
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
