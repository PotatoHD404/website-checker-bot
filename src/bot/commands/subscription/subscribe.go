package subscription

import (
	. "website-checker-bot/bot/commands/env"

	"gopkg.in/telebot.v3"
)

func HandleSubscribe(env *Env, c telebot.Context, args []string) error {
	if len(args) != 1 {
		err := c.Reply("Usage: /subscribe <name>")
		if err != nil {
			return err
		}
		return nil
	}

	websiteName := args[0]
	if !env.Db.CheckWebsite(websiteName) {
		err := c.Reply("Website with this name does not exist")
		if err != nil {
			return err
		}
		return nil
	}

	userId := c.Sender().ID
	if env.Db.CheckSubscription(userId, websiteName) {
		err := c.Reply("You already subscribed to this website")
		if err != nil {
			return err
		}
		return nil
	}

	env.Db.AddSubscription(userId, websiteName)

	err := c.Reply("Subscribed")

	if err != nil {
		return err
	}

	return nil
}
