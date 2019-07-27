package cmd

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"

	"github.com/urfave/cli"
)

var getWebhookInfoCommand = cli.Command{
	Name:     "get-webhook-info",
	Aliases:  []string{"getWebhookInfo"},
	Category: "webhook",
	Usage:    "get current webhook status.",
	Action: internal.NewInfoAction(func(ctx context.Context, _ *cli.Context, client *tg.Client) (interface{}, error) {
		return client.GetWebhookInfo(ctx)
	}),
	Flags: flags(),
}
