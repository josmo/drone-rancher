package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

type Rancher struct {
	Url            string `json:"url"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	Service        string `json:"service"`
	Image          string `json:"docker_image"`
	StartFirst     bool   `json:"start_first"`
	Confirm        bool   `json:"confirm"`
	Timeout        int    `json:"timeout"`
	IntervalMillis int64  `json:"interval_millis"`
	BatchSize      int64  `json:"batch_size"`
}

var build string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "rancher publish"
	app.Usage = "rancher publish"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.0+%s", build)

	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "url",
			Usage:  "url to the rancher api",
			EnvVar: "PLUGIN_URL",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "rancher access key",
			EnvVar: "PLUGIN_ACCESS_KEY, RANCHER_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "rancher secret key",
			EnvVar: "PLUGIN_SECRET_KEY, RANCHER_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "service",
			Usage:  "Service to act on",
			EnvVar: "PLUGIN_SERVICE",
		},
		cli.StringFlag{
			Name:   "docker-image",
			Usage:  "image to use",
			EnvVar: "PLUGIN_DOCKER_IMAGE",
		},
		cli.BoolTFlag{
			Name:   "start-first",
			Usage:  "Start new container before stoping old",
			EnvVar: "PLUGIN_START_FIRST",
		},
		cli.BoolFlag{
			Name:   "confirm",
			Usage:  "auto confirm the service upgrade if successful",
			EnvVar: "PLUGIN_CONFIRM",
		},
		cli.IntFlag{
			Name:   "timeout",
			Usage:  "the maximum wait time in seconds for the service to upgrade",
			Value:  30,
			EnvVar: "PLUGIN_TIMEOUT",
		},
		cli.Int64Flag{
			Name:   "interval-millis",
			Usage:  "The interval for batch size upgrade",
			Value:  1000,
			EnvVar: "PLUGIN_INTERVAL_MILLIS",
		},
		cli.Int64Flag{
			Name:   "batch-size",
			Usage:  "The upgrade batch size",
			Value:  1,
			EnvVar: "PLUGIN_BATCH_SIZE",
		},
		cli.BoolTFlag{
			Name:   "yaml-verified",
			Usage:  "Ensure the yaml was signed",
			EnvVar: "DRONE_YAML_VERIFIED",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		URL:            c.String("url"),
		Key:            c.String("access-key"),
		Secret:         c.String("secret-key"),
		Service:        c.String("service"),
		DockerImage:    c.String("docker-image"),
		StartFirst:     c.BoolT("start-first"),
		Confirm:        c.Bool("confirm"),
		Timeout:        c.Int("timeout"),
		IntervalMillis: c.Int64("interval-millis"),
		BatchSize:      c.Int64("batch-size"),
		YamlVerified:   c.BoolT("yaml-verified"),
	}
	return plugin.Exec()
}
