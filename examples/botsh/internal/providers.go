package internal

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
)

func provideOutput(cliCtx *cli.Context) Output {
	return NewOutput(
		cliCtx.App.Writer,
		cliCtx.String("format"),
	)
}

func provideClient(cliCtx *cli.Context) *tg.Client {
	domain := cliCtx.GlobalString("api-domain")

	transport := tg.NewHTTPTransport(
		tg.WithHTTPBuildCallURLFunc(func(token, method string) string {
			return fmt.Sprintf("https://%s/bot%s/%s", domain, token, method)
		}),
		tg.WithHTTPBuildFileURLFunc(func(token, path string) string {
			return fmt.Sprintf("https://%s/file/bot%s/%s", domain, token, path)
		}),
	)

	return tg.NewClient(cliCtx.GlobalString("token"),
		tg.WithTransport(transport),
	)
}

func provideCtx(cliCtx *cli.Context) (context.Context, context.CancelFunc) {
	timeout := cliCtx.GlobalDuration("request-timeout")

	return context.WithTimeout(context.Background(), timeout)
}
