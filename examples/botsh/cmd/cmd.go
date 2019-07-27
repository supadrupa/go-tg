package cmd

import (
	"time"

	"github.com/urfave/cli"
)

func flags(fs ...cli.Flag) []cli.Flag {
	return append(fs, cli.StringFlag{
		Name:  "format, f",
		Usage: "Output format (json, pretty or custom template)",
		Value: "pretty",
	})
}

func NewApp() *cli.App {
	app := cli.NewApp()

	app.Name = "botsh"
	app.Authors = []cli.Author{
		{
			Name:  "Sasha Savchuk",
			Email: "mrxlinch@gmail.com",
		},
	}
	app.EnableBashCompletion = true
	app.Usage = "Simple Telegram Bot API command-line client"
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			Usage:  "Telegram Bot API token",
			EnvVar: "TELEGRAM_BOT_TOKEN",
		},
		cli.DurationFlag{
			Name:  "request-timeout",
			Usage: "timeout for requests",
			Value: time.Minute,
		},
		cli.StringFlag{
			Name:   "api-domain",
			Usage:  "Telegram Bot API domain",
			EnvVar: "TELEGRAM_BOT_API_DOMAIN",
			Value:  "api.telegram.org",
		},
	}

	app.Commands = []cli.Command{
		getMeCommand,
		getWebhookInfoCommand,
		setWebhookCommand,
		getChatCommand,
		getFileCommand,
		getChatAdmins,
		getChatMembersCount,
		deleteWebhookCommand,
		getUpdatesCommand,
	}

	return app
}
