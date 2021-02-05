package main

import "github.com/urfave/cli"

const (
	OptProject   = "project"
	OptOutputDir = "out"
)

var flags = []cli.Flag{
	cli.StringSliceFlag{
		Name:  OptProject + ", p",
		Usage: "go project",
	},
	cli.StringFlag{
		Name:  OptOutputDir + ", o",
		Usage: "output directory which aggregate all projects' vendor source",
	},
}
