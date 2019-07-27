package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"
)

var setWebhookCommand = cli.Command{
	Name:      "set-webhook",
	Aliases:   []string{"setWebhook"},
	Category:  "webhook",
	Usage:     "use this method to specify a url and receive incoming updates via an outgoing webhook.",
	ArgsUsage: "URL",

	Before: func(cliCtx *cli.Context) error {
		validateURL := func() error {
			_, err := url.ParseRequestURI(cliCtx.Args().First())
			if err != nil {
				return fmt.Errorf("url: Invalid or missing")
			}

			return nil
		}

		validateAllowedUpdates := func() error {
			uts := cliCtx.StringSlice("allowed-updates")

			errs := []error{}

			for _, v := range uts {
				if !isValidUpdateType(v) {
					errs = append(errs, fmt.Errorf("--allowed-updates: invalid update type: '%s'", v))
				}
			}

			if len(errs) > 0 {
				return cli.NewMultiError(errs...)
			}

			return nil
		}

		validateCertificate := func() error {
			if cliCtx.IsSet("certificate-path") {
				return isFileAvailable(cliCtx.String("certificate-path"))
			}

			return nil
		}

		return validate(
			validateURL,
			validateAllowedUpdates,
			validateCertificate,
		)
	},

	Action: internal.NewAction(func(ctx context.Context, cliCtx internal.CLIContext, client *tg.Client, output internal.Output) error {
		url := cliCtx.Args().First()
		allowedUpdates, _ := parseAllowedUpdates(cliCtx.StringSlice("allowed-updates"))
		maxConnections := cliCtx.Int("max-connections")

		if !cliCtx.IsSet("allowed-updates") {
			allowedUpdates = nil
		}

		err := client.SetWebhook(ctx, url, &tg.WebhookOptions{
			MaxConnections: maxConnections,
			AllowedUpdates: allowedUpdates,
		})

		return err
	}),

	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "max-connections, m",
			Usage: "maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery, 1-100.",
			Value: 40,
		},
		cli.StringFlag{
			Name:  "certificate-path, c",
			Usage: "load public key certificate from `FILE`",
		},
		cli.StringSliceFlag{
			Name: "allowed-updates, u",
			Usage: fmt.Sprintf(
				"list the types of updates you want your bot to receive (values: %s)",
				strings.Join(updateTypes, ", "),
			),
		},
	},
}
