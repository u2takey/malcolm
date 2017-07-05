package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	router "github.com/u2takey/malcolm/server"
	"github.com/urfave/cli"
)

// Command exports the server command.
var Command = cli.Command{
	Name:   "server",
	Usage:  "starts the malcolm server daemon",
	Action: server,
	Flags: []cli.Flag{
		cli.BoolFlag{
			EnvVar: "DEBUG",
			Name:   "debug",
			Usage:  "start the server in debug mode",
		},
		cli.StringFlag{
			EnvVar: "SERVER_HOST",
			Name:   "server-host",
			Usage:  "server host",
		},
		cli.StringFlag{
			EnvVar: "SERVER_ADDR",
			Name:   "server-addr",
			Usage:  "server address",
			Value:  ":7700",
		},
	},
}

func server(c *cli.Context) error {
	// debug level if requested by user
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	// setup the server and start the listener
	handler := router.Load(
		ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true),
	)

	// start the server
	return http.ListenAndServe(
		c.String("server-addr"),
		handler,
	)

	return nil
}
