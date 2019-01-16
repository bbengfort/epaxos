package main

import (
	"os"
	"encoding/json"
	"io/ioutil"

	"github.com/bbengfort/epaxos"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// Comand variables
var (
	config *epaxos.Config
	replica *epaxos.Replica
	// client * epaxos.Client
)

func main() {
	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "epaxos"
	app.Version = epaxos.PackageVersion
	app.Usage = "implements the ePaxos consensus algorithm"
	app.Before = initConfig

	// Global Arguments
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c, config",
			Usage: "configuration file for replica",
		},
	}

	// Define the commands available to the application
	app.Commands = []cli.Command{
		{
			Name:     "serve",
			Usage:    "run an epaxos server",
			Action:   serve,
			Category: "server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "n, name",
					Usage: "unique name of the replica instance",
				},
				cli.DurationFlag{
					Name:  "u, uptime",
					Usage: "specify a duration for the server to run",
				},
				cli.StringFlag{
					Name:  "o, outpath",
					Usage: "write metrics out to the specified path",
				},
				cli.Int64Flag{
					Name:  "s, seed",
					Usage: "specify the random seed",
				},
			},
		},
		{
			Name:     "commit",
			Usage:    "commit an entry to the distributed log",
			Action:   commit,
			Category: "client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "name or address of replica to connect to",
				},
				cli.StringFlag{
					Name:  "k, key",
					Usage: "the key to associate the entry with",
				},
				cli.StringFlag{
					Name:  "v, value",
					Usage: "the value to commit with the associated key",
				},
			},
		},
		{
			Name:     "bench",
			Usage:    "run an epaxos benchmark with concurrent network",
			Action:   bench,
			Category: "client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "name or address of replica to connect to",
					Value: "",
				},
				cli.UintFlag{
					Name:  "r, requests",
					Usage: "number of requests issued per client",
					Value: 1000,
				},
				cli.IntFlag{
					Name:  "s, size",
					Usage: "number of bytes per value",
					Value: 32,
				},
				cli.DurationFlag{
					Name:  "d, delay",
					Usage: "wait specified time before starting benchmark",
				},
				cli.IntFlag{
					Name:  "i, indent",
					Usage: "indent the results by specified number of spaces",
				},
				cli.BoolFlag{
					Name:  "b, blast",
					Usage: "send all requests per client at once",
				},
			},
		},
	}

	// Run the CLI program
	app.Run(os.Args)
}

//===========================================================================
// Initialization
//===========================================================================

func initConfig(c *cli.Context) (err error) {
	config = new(epaxos.Config)
	if cpath := c.String("config"); cpath != "" {
		var data []byte
		if data, err = ioutil.ReadFile(cpath); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err = json.Unmarshal(data, &config); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

//===========================================================================
// Server Commands
//===========================================================================

func serve(c *cli.Context) (err error) {
	if name := c.String("name"); name != "" {
		config.Name = name
	}

	if seed := c.Int64("seed"); seed > 0 {
		config.Seed = seed
	}

	if replica, err = epaxos.New(config); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = replica.Listen(); err != nil{
		return cli.NewExitError(err, 1)
	}

	return nil
}

//===========================================================================
// Client Commands
//===========================================================================

func commit(c *cli.Context) (err error) {
	return cli.NewExitError("not implemented yet", 1)
}

func bench(c *cli.Context) (err error) {
	return cli.NewExitError("not implemented yet", 1)
}
