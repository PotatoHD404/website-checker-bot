package admin

import (
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/database/models"
	"website-checker-bot/utils"

	"strconv"

	"gopkg.in/telebot.v3"
)

func HandleAddAdmin(env *Env, c telebot.Context, args []string) error {
	admins := env.Db.GetAdmins()

	if len(admins) == 0 {
		env.Db.AddAdmin(c.Sender().ID)
		err := c.Reply("You are admin now")
		if err != nil {
			return err
		}
		return nil
	} else if !utils.Contains(env.Db.GetAdmins(), NewAdmin(c.Sender().ID)) {
		c.Reply("You are not admin")
		return nil
	}

	if len(args) == 0 {
		c.Reply("You are already admin")
		return nil
	}

	if len(args) > 1 {
		c.Reply("Too many arguments")
		return nil
	}

	userId, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		c.Reply("Invalid argument")
		return nil
	}

	if utils.Contains(admins, NewAdmin(userId)) {
		c.Reply("User is already admin")
	}

	env.Db.AddAdmin(userId)
	err = c.Reply("User is admin now")
	if err != nil {
		return err
	}

	return nil
}
