package cmd

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"

	"github.com/urfave/cli"
)

var getMeCommand = cli.Command{
	Name:     "get-me",
	Aliases:  []string{"getMe"},
	Category: "generic",
	Usage:    "returns basic information about the bot.",
	Action: internal.NewInfoAction(func(ctx context.Context, _ internal.CLIContext, client *tg.Client) (interface{}, error) {
		return client.GetMe(ctx)
	}),
	Flags: flags(),
}
