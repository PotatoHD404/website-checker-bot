package commands

import (
	"gopkg.in/telebot.v3"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/middlewares"
)

func HandleHelp(env *Env, c telebot.Context, args []string) error {
	msg := "Available commands:\n" +
		"/start - start the bot\n" +
		"/help - show this message\n" +
		//"/addwebsite <url> <name> - add a website to the database\n" +
		//"/deletewebsite <name> - delete a website from the database\n" +
		"/get_websites - get all websites from the database\n" +
		"/subscribe <name> - subscribe to a website\n" +
		"/unsubscribe <name> - unsubscribe from a website\n" +
		"/get_subscriptions - get all your subscriptions\n"
	//"/addadmin <id> - add a user to the admin list\n" +
	//"/deleteadmin <id> - delete a user from the admin list\n" +
	//"/getadmins - get all admins\n"
	if CheckAdmin(env, c) {
		msg += "Admin commands:\n" +
			"/add_website <url> <name> - add a website to the database\n" +
			"/delete_website <name> - delete a website from the database\n" +
			"/add_admin <id> - add a user to the admin list\n" +
			"/delete_admin <id> - delete a user from the admin list\n" +
			"/get_admins - get all admins\n"
	}

	err := c.Reply(msg)
	if err != nil {
		return err
	}

	return nil
}
