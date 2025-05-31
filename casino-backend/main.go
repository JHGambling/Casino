package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	// For more info on creating the CLI:
	// https://cli.urfave.org/v3/getting-started/
	cmd := &cli.Command{
		Name:                  "casino",
		Usage:                 "CLI",
		EnableShellCompletion: true,

		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Starts a casino server instance",
				Action: func(c context.Context, cmd *cli.Command) error {
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
