package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"
)

var getUpdatesCommand = cli.Command{
	Name:     "get-updates",
	Aliases:  []string{"getMe"},
	Category: "generic",
	Usage:    "returns basic information about the bot.",

	Action: internal.NewAction(func(
		ctx context.Context,
		cliCtx internal.CLIContext,
		client *tg.Client,
		output internal.Output,
	) error {
		ctx = context.Background()

		updates := make(chan tg.Update, cliCtx.Int("limit"))

		g, ctx := errgroup.WithContext(ctx)

		// writer
		g.Go(func() error {
			offset := tg.UpdateID(0)
			limit := cliCtx.Int("limit")
			timeout := cliCtx.Duration("timeout")

			for {
				upds, err := client.GetUpdates(ctx, &tg.UpdatesOptions{
					Offset:  offset,
					Limit:   limit,
					Timeout: timeout,
				})

				if err != nil {
					fmt.Fprintf(output, "get updates error: %v\n", err)
					time.Sleep(time.Second * 5)
					continue
				}

				if len(upds) > 0 {
					for _, update := range upds {
						updates <- update
					}

					offset = upds[len(upds)-1].ID + 1
				}
			}

			return nil
		})

		// reader
		g.Go(func() error {
			for update := range updates {
				var body interface{}

				switch update.Type() {
				case tg.UpdateMessage:
					body = update.Message
				case tg.UpdateEditedMessage:
					body = update.EditedMessage
				case tg.UpdateChannelPost:
					body = update.ChannelPost
				case tg.UpdateEditedChannelPost:
					body = update.EditedChannelPost
				case tg.UpdateInlineQuery:
					body = update.InlineQuery
				case tg.UpdateChosenInlineResult:
					body = update.ChosenInlineResult
				case tg.UpdateCallbackQuery:
					body = update.CallbackQuery
				case tg.UpdateShippingQuery:
					body = update.ShippingQuery
				case tg.UpdatePreCheckoutQuery:
					body = update.PreCheckoutQuery
				case tg.UpdatePoll:
					body = update.Poll
				default:
					body = update
				}

				output.Print(body)
			}

			return nil
		})

		return g.Wait()
	}),

	Flags: flags(
		cli.IntFlag{
			Name:  "limit",
			Usage: "limits the number of updates to be retrieved once",
			Value: 100,
		},
		cli.DurationFlag{
			Name:  "timeout",
			Usage: "timeout for long polling (1s-100s)",
			Value: time.Second * 10,
		},
		cli.StringSliceFlag{
			Name: "allowed-updates, u",
			Usage: fmt.Sprintf(
				"list the types of updates you want your bot to receive (values: %s)",
				strings.Join(updateTypes, ", "),
			),
		},
	),
}
