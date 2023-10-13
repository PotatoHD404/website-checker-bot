package bot

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	. "website-checker-bot/bot/commands"
	. "website-checker-bot/bot/commands/admin"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/commands/subscription"
	. "website-checker-bot/bot/commands/website"

	"gopkg.in/telebot.v3"
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

func (b *Bot) Init(e *Env) {
	b.bot.Handle(telebot.OnText, func(c telebot.Context) error {
		r, _ := regexp.Compile(`^\/([a-z_]+)( [a-zA-ZА-Яа-я0-9_\/:.\-=?&]+)*$`)
		message := c.Text()
		if !r.MatchString(message) {
			err := c.Reply("I don't understand you. Use /help to see all available commands.")
			if err != nil {
				fmt.Println(err)
				panic("can't send message")
			}
			return nil
		}
		submatch := strings.Split(message, " ")
		command := submatch[0]
		args := submatch[1:]
		switch command {
		case "/start":
			return HandleStart(e, c, args)
		case "/help":
			return HandleHelp(e, c, args)
		case "/add_website":
			return HandleAddWebsite(e, c, args)
		case "/get_websites":
			return HandleGetWebsites(e, c, args)
		case "/delete_website":
			return HandleDeleteWebsite(e, c, args)
		case "/add_admin":
			return HandleAddAdmin(e, c, args)
		case "/get_admins":
			return HandleGetAdmins(e, c, args)
		case "/delete_admin":
			return HandleDeleteAdmin(e, c, args)
		case "/subscribe":
			return HandleSubscribe(e, c, args)
		case "/unsubscribe":
			return HandleUnsubscribe(e, c, args)
		case "/get_subscriptions":
			return HandleGetSubscriptions(e, c, args)

		default:
			err := c.Reply("There is no such command. Use /help to see all available commands.")
			if err != nil {
				fmt.Println(err)
				panic("can't send message")
			}
			return nil
		}
	})

	b.bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
		err := c.Reply("I can't handle photos yet")
		return err
	})

	b.bot.Handle(telebot.OnDocument, func(c telebot.Context) error {
		err := c.Reply("I can't handle documents yet")
		return err
	})

	b.bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		err := c.Reply("I can't handle stickers yet")
		return err
	})

	b.bot.Handle(telebot.OnAudio, func(c telebot.Context) error {
		err := c.Reply("I can't handle audio yet")
		return err
	})

	b.bot.Handle(telebot.OnVoice, func(c telebot.Context) error {
		err := c.Reply("I can't handle voice yet")
		return err
	})

	b.bot.Handle(telebot.OnVideo, func(c telebot.Context) error {
		err := c.Reply("I can't handle video yet")
		return err
	})

	b.bot.Handle(telebot.OnVideoNote, func(c telebot.Context) error {
		err := c.Reply("I can't handle video notes yet")
		return err
	})

	b.bot.Handle(telebot.OnContact, func(c telebot.Context) error {
		err := c.Reply("I can't handle contacts yet")
		return err
	})

	b.bot.Handle(telebot.OnLocation, func(c telebot.Context) error {
		err := c.Reply("I can't handle locations yet")
		return err
	})

	b.bot.Handle(telebot.OnVenue, func(c telebot.Context) error {
		err := c.Reply("I can't handle venues yet")
		return err
	})

	b.bot.Handle(telebot.OnPoll, func(c telebot.Context) error {
		err := c.Reply("I can't handle polls yet")
		return err
	})

	b.bot.Handle(telebot.OnDice, func(c telebot.Context) error {
		err := c.Reply("I can't handle dice yet")
		return err
	})

	b.bot.Handle(telebot.OnGame, func(c telebot.Context) error {
		err := c.Reply("I can't handle games yet")
		return err
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
