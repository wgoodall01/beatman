package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "beatman"
	app.Usage = "Manage a Beat Saber library"

	app.Commands = []cli.Command{cliServe()}

	if err := app.Run(os.Args); err != nil {
		log.cli.Fatal(err)
	}
}

func cliServe() cli.Command {
	c := cli.Command{}
	c.Name = "serve"
	c.Usage = "Serve the webUI on the given port"

	c.Flags = []cli.Flag{
		cli.StringFlag{Name: "port,p", EnvVar: "PORT", Value: "8080", Usage: "Port to listen on"},
		cli.IntFlag{Name: "sync-rate", Value: 30, Usage: "Library refresh rate, in seconds"},
		cli.StringFlag{Name: "library,l", Value: "_beatsaber", Usage: "Path to BeatSaber directory"},
	}

	c.Action = func(c *cli.Context) error {
		// Set up beatsaber library
		library, err := NewLibrary(c.String("library"))
		if err != nil {
			log.cli.WithError(err).Fatal("Could not open library")
		}
		library.StartSync(time.Duration(c.Int64("sync-rate")) * time.Second)

		// Start serving the API
		api := &API{library, ":" + c.String("port")}
		api.Serve()

		return nil
	}

	return c
}
