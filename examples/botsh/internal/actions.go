package internal

import (
	"context"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
)

type CLIContext = *cli.Context

func NewAction(do func(ctx context.Context, cliCtx CLIContext, client *tg.Client, output Output) error) cli.ActionFunc {
	return cli.ActionFunc(func(cliCtx *cli.Context) error {
		ctx, cancel := provideCtx(cliCtx)
		defer cancel()

		client := provideClient(cliCtx)
		output := provideOutput(cliCtx)

		return do(ctx, cliCtx, client, output)
	})
}

func NewInfoAction(do func(ctx context.Context, cliCtx CLIContext, client *tg.Client) (interface{}, error)) cli.ActionFunc {
	return NewAction(func(ctx context.Context, cliCtx CLIContext, client *tg.Client, output Output) error {
		object, err := do(ctx, cliCtx, client)
		if err != nil {
			return err
		}
		return output.Print(object)
	})
}
