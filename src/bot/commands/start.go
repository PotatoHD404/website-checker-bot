package commands

import (
	. "website-checker-bot/bot/commands/env"

	"gopkg.in/telebot.v3"
)

func HandleStart(env *Env, c telebot.Context, args []string) error {
	err := c.Reply("Hello! I'm a bot that will notify you when a website is changed.\nUse /help to see all available commands.")
	if err != nil {
		return err
	}
	return nil
}
