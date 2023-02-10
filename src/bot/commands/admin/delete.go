package admin

import (
	"strconv"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/middlewares"
	. "website-checker-bot/database/models"
	"website-checker-bot/utils"

	"gopkg.in/telebot.v3"
)

func HandleDeleteAdmin(env *Env, c telebot.Context, args []string) error {
	if !CheckAdmin(env, c) {
		return nil
	}

	if len(args) != 1 {
		err := c.Reply("Usage: /delete_admin <userId>")
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

	if !utils.Contains(env.Db.GetAdmins(), NewAdmin(userId, username)) {
		err := c.Reply("User is not admin")
		if err != nil {
			return err
		}
		return nil
	}

	env.Db.DeleteAdmin(userId)
	err = c.Reply("Admin deleted successfully")
	if err != nil {
		return err
	}

	return nil
}
