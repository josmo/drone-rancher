package main

import (
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/joho/godotenv/autoload"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "Drone Rancher"
	app.Usage = "Drone Rancher usage"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{

		//
		// plugin args
		//
		cli.StringFlag{
			Name:   "rancher.url",
			Usage:  "rancher url",
			EnvVar: "PLUGIN_URL",
		},
		cli.StringFlag{
			Name:   "rancher.accesskey",
			Usage:  "rancher accesskey",
			EnvVar: "PLUGIN_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "rancher.secretkey",
			Usage:  "rancher secretkey",
			EnvVar: "PLUGIN_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "rancher.service",
			Usage:  "rancher service",
			EnvVar: "PLUGIN_SERVICE",
		},
		cli.StringFlag{
			Name:   "rancher.image",
			Usage:  "rancher image",
			EnvVar: "PLUGIN_IMAGE",
		},
		cli.BoolFlag{
			Name:   "rancher.startfirst",
			Usage:  "rancher startfirst",
			EnvVar: "PLUGIN_START_FIRST",
		},
		cli.BoolFlag{
			Name:   "rancher.confirm",
			Usage:  "rancher confirm",
			EnvVar: "PLUGIN_CONFIRM",
		},
		cli.IntFlag{
			Name:   "rancher.timeout",
			Usage:  "rancher timeout",
			EnvVar: "PLUGIN_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "rancher.params",
			Usage:  "rancher params",
			EnvVar: "PLUGIN_PARAMS",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Rancher: Rancher{
			Url:        c.String("rancher.url"),
			AccessKey:  c.String("rancher.accesskey"),
			SecretKey:  c.String("rancher.secretkey"),
			Service:    c.String("rancher.service"),
			Image:      c.String("rancher.image"),
			StartFirst: c.Bool("rancher.startfirst"),
			Confirm:    c.Bool("rancher.confirm"),
			Timeout:    int64(c.Int("rancher.timeout")),
			Params:     c.String("rancher.params"),
		},
	}

	if err := plugin.Exec(); err != nil {
		return err
	}

	return nil
}
