package website

import (
	"regexp"
	. "website-checker-bot/bot/commands/env"
	. "website-checker-bot/bot/middlewares"

	"gopkg.in/telebot.v3"
)

func HandleAddWebsite(env *Env, c telebot.Context, args []string) error {
	if !CheckAdmin(env, c) {
		return nil
	}

	if len(args) < 2 || len(args) > 3 {
		err := c.Reply("Usage: /add_website <url> <name> [xpath]")
		if err != nil {
			return err
		}
		return nil
	}

	websiteUrl := args[0]
	websiteName := args[1]
	var websiteXpath string
	if len(args) == 3 {
		websiteXpath = args[2]
	} else {
		websiteXpath = ""
	}
	if env.Db.CheckWebsite(websiteName) {
		err := c.Reply("Website with this name already exists")
		if err != nil {
			return err
		}
		return nil
	}

	if env.Db.CheckWebsiteUrl(websiteUrl) {
		err := c.Reply("Website with this url already exists")
		if err != nil {
			return err
		}
		return nil
	}

	r := regexp.MustCompile(`^(http|https)://[a-z0-9]+([\-\.][a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	if !r.MatchString(websiteUrl) {
		err := c.Reply("Invalid URL")
		if err != nil {
			return err
		}
		return nil
	}
	r = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !r.MatchString(websiteName) {
		err := c.Reply("Invalid name")
		if err != nil {
			return err
		}
		return nil
	}
	// regexp for xpath
	r = regexp.MustCompile(`^$|^(//|/)[a-zA-Z0-9]+(//|/)[a-zA-Z0-9]+$`)
	if !r.MatchString(websiteXpath) {
		err := c.Reply("Invalid xpath")
		if err != nil {
			return err
		}
		return nil
	}

	env.Db.AddWebsite(websiteName, websiteUrl, websiteXpath)
	env.Db.AddSubscription(c.Sender().ID, websiteName)

	err := c.Reply("Website added")
	if err != nil {
		return err
	}

	return nil
}
