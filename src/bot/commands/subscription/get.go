package subscription

import (
	. "website-checker-bot/bot/commands/env"

	"gopkg.in/telebot.v3"
)

func HandleGetSubscriptions(env *Env, c telebot.Context, args []string) error {
	subscriptions := env.Db.GetSubscriptions(c.Sender().ID)

	message := "Subscriptions:\n"
	for _, subscription := range subscriptions {
		message += subscription.Name + " - " + subscription.Url + "\n"
	}
	err := c.Reply(message)
	if err != nil {
		return err
	}

	return nil
}
