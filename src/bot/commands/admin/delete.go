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

	if len(args) == 0 {
		c.Send("Please provide an admin ID")
		return nil
	}

	if len(args) > 1 {
		c.Send("Please provide only one admin ID")
		return nil
	}

	userId, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		c.Send("Invalid argument")
		return nil
	}

	if !utils.Contains(env.Db.GetAdmins(), NewAdmin(userId)) {
		c.Send("User is not admin")
		return nil
	}

	env.Db.DeleteAdmin(userId)

	return nil
}
