package website

import (
	. "website-checker-bot/bot/commands/env"

	"gopkg.in/telebot.v3"
)

func HandleGetWebsites(env *Env, c telebot.Context, args []string) error {

	websites := env.Db.GetWebsites(false)
	message := "Websites:\n"
	for _, website := range websites {
		message += website.Name + " - " + website.Url + "\n"
	}
	err := c.Reply(message)
	if err != nil {
		return err
	}

	return nil
}
