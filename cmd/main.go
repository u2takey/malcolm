package main

import (
	"fmt"
	"os"

	"github.com/u2takey/malcolm/cmd/server"
	"github.com/u2takey/malcolm/version"

	"github.com/ianschenck/envflag"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

func main() {
	envflag.Parse()

	app := cli.NewApp()
	app.Name = "malcolm"
	app.Version = version.Version.String()
	app.Usage = "command line utility"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "s, server",
			Usage:  "server location",
			EnvVar: "MALCOLM_SERVER",
		},
	}
	app.Commands = []cli.Command{
		server.Command,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
