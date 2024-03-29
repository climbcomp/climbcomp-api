package cmd

import (
	"github.com/climbcomp/climbcomp-go/climbcomp"
	"github.com/urfave/cli"
)

// NewApp returns the App
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Version = climbcomp.VERSION
	app.Usage = "A competition climbing API"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:   "server",
			Usage:  "Starts the climbcomp server",
			Action: OnServerCmd,
		},
		{
			Name:  "meta",
			Usage: "Meta API commands",
			Subcommands: []cli.Command{
				{
					Name:   "version",
					Usage:  "Returns the server version",
					Action: OnMetaVersionCmd,
				},
			},
		},
	}

	return app
}
