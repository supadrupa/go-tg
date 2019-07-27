package cmd

import (
	"context"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"
)

var getChatCommand = cli.Command{
	Name:      "get-chat",
	Aliases:   []string{"getChat"},
	Category:  "chats",
	Usage:     "get information about chat",
	ArgsUsage: "PEER_ID",

	Action: internal.NewInfoAction(func(ctx context.Context, cliCtx internal.CLIContext, client *tg.Client) (chat interface{}, err error) {
		peer, err := tg.ParsePeer(cliCtx.Args().First())
		if err != nil {
			return nil, err
		}

		return client.GetChat(ctx, peer)
	}),

	Flags: flags(),
}
