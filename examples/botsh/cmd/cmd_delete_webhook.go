package cmd

import (
	"context"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"
)

var deleteWebhookCommand = cli.Command{
	Name:     "delete-webhook",
	Aliases:  []string{"deleteWebhook"},
	Category: "webhook",
	Usage:    "removes current set webhook",

	Action: internal.NewAction(func(ctx context.Context, cliCtx internal.CLIContext, client *tg.Client, output internal.Output) error {
		return client.DeleteWebhook(ctx)
	}),

	Flags: flags(),
}
