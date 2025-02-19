package main

import (
	"fmt"
	"github.com/jgrigorian/certscan/cmd/list"
	"github.com/jgrigorian/certscan/cmd/show"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func main() {
	//list.Certificates()
	// CLI
	app := &cli.App{
		Name:  "certscan",
		Usage: "Tool to interact with tls secrets in a kubernetes cluster",
		OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(ctx.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
		Commands: []*cli.Command{
			///////////////////////////////////////////////////////////
			// LIST
			///////////////////////////////////////////////////////////
			{
				Name:  "list",
				Usage: "Option for listing certificates",
				Subcommands: []*cli.Command{
					{
						Name:  "certificates",
						Usage: "list certificates",
						Action: func(c *cli.Context) error {
							list.Certificates(c)
							return nil
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "expiring",
								Aliases: []string{"e"},
								Usage:   "List only expired or expiring certificates",
							},
							&cli.StringFlag{
								Name:    "namespace",
								Aliases: []string{"n"},
								Usage:   "Desired namespace (Example: default)",
							},
							&cli.BoolFlag{
								Name:    "all-namespaces",
								Aliases: []string{"A"},
								Usage:   "All namespaces",
							},
						},
					},
				},
			},
			///////////////////////////////////////////////////////////
			// SHOW
			///////////////////////////////////////////////////////////
			{
				Name:  "show",
				Usage: "Option for showing certificate details",
				Subcommands: []*cli.Command{
					{
						Name:  "certificate",
						Usage: "show certificate",
						Action: func(c *cli.Context) error {
							show.Certificate(c)
							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "secret",
								Aliases: []string{"s"},
								Usage:   "Name of the desired secret. Example: staging-ssl-secret",
							},
							&cli.StringFlag{
								Name:    "namespace",
								Aliases: []string{"n"},
								Usage:   "Desired namespace (Example: default). NOTE: If no namespace is provided, default is used.",
							},
						},
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
