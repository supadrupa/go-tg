package cmd

import (
	"context"
	"io"
	"os"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/examples/botsh/internal"
)

var getFileCommand = cli.Command{
	Name:      "get-file",
	Aliases:   []string{"getFile"},
	Category:  "generic",
	Usage:     "get information about chat",
	ArgsUsage: "FILE_ID [PATH]",

	Action: internal.NewAction(func(ctx context.Context, cliCtx internal.CLIContext, client *tg.Client, output internal.Output) error {
		file, err := client.GetFile(ctx, tg.FileID(cliCtx.Args().First()))

		if err != nil {
			return err
		}

		// download file
		if cliCtx.NArg() > 1 {
			var downloadTo io.Writer

			pathOpt := cliCtx.Args().Get(1)

			// write file to std output
			if pathOpt == "-" {
				downloadTo = output
			} else {
				outputFile, err := os.Create(pathOpt)
				if err != nil {
					return err
				}
				defer outputFile.Close()

				downloadTo = outputFile
			}

			reader, err := file.NewReader(ctx)
			if err != nil {
				return err
			}
			defer reader.Close()

			_, err = io.Copy(downloadTo, reader)
			return err
		}

		return output.Print(file)
	}),

	Flags: flags(),
}
