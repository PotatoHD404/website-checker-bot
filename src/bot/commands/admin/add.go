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
		env.Db.AddAdmin(c.Sender().ID, c.Sender().Username)
		err := c.Reply("You are admin now")
		if err != nil {
			return err
		}
		return nil
	} else if !utils.Contains(env.Db.GetAdmins(), NewAdmin(c.Sender().ID, c.Sender().Username)) {
		err := c.Reply("You are not admin")
		if err != nil {
			return err
		}
		return nil
	}

	if len(args) == 0 {
		err := c.Reply("You are already admin")
		if err != nil {
			return err
		}
		return nil
	}

	if len(args) != 1 {
		err := c.Reply("Usage: /addadmin <user_id>")
		if err != nil {
			return err
		}
		return nil
	}

	userId, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		err := c.Reply("Invalid argument")
		if err != nil {
			return err
		}
		return nil
	}

	username := c.Sender().Username

	if utils.Contains(admins, NewAdmin(userId, username)) {
		err := c.Reply("User is already admin")
		if err != nil {
			return err
		}
		return nil
	}

	env.Db.AddAdmin(userId, username)
	err = c.Reply("User is admin now")
	if err != nil {
		return err
	}

	return nil
}
