package server

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/arlert/malcolm/model"
	router "github.com/arlert/malcolm/server"
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
			EnvVar: "SERVER_ADDR",
			Name:   "server-addr",
			Usage:  "server address",
			Value:  ":7700",
		},
		cli.StringFlag{
			EnvVar: "MONGO_HOST",
			Name:   "mongo-host",
			Usage:  "mogno host",
			Value:  "mongo:27017",
		},
		cli.StringFlag{
			EnvVar: "MONGO_DB",
			Name:   "mongo-db",
			Usage:  "mogno db",
			Value:  "malcolm",
		},
		cli.BoolFlag{
			EnvVar: "MONGO_DIRECT",
			Name:   "mongo-direct",
			Usage:  "mongo direct",
		},
		cli.StringFlag{
			EnvVar: "KUBERNETE_ADDR",
			Name:   "kubernete-addr",
			Usage:  "kubernete addr",
			Value:  "https://kubernetes.default",
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
	cfg := &model.Config{}
	cfg.MgoCfg.Host = c.String("mongo-host")
	cfg.MgoCfg.DB = c.String("mongo-db")
	cfg.MgoCfg.Direct = c.Bool("mongo-direct")
	cfg.KubeAddr = c.String("kubernete-addr")
	// setup the server and start the listener
	handler := router.Load(cfg)

	// start the server
	return http.ListenAndServe(
		c.String("server-addr"),
		handler,
	)

	return nil
}
