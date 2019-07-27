package cmd

import (
	"context"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"
)

var getChatAdmins = cli.Command{
	Name:      "get-chat-administrators",
	Aliases:   []string{"getChatAdministrators"},
	Category:  "chats",
	Usage:     "get a list of administrators in a chat",
	ArgsUsage: "PEER_ID",

	Action: internal.NewInfoAction(func(ctx context.Context, cliCtx internal.CLIContext, client *tg.Client) (chat interface{}, err error) {
		peer, err := tg.ParsePeer(cliCtx.Args().First())
		if err != nil {
			return nil, err
		}

		return client.GetChatAdministrators(ctx, peer)
	}),

	Flags: flags(),
}
