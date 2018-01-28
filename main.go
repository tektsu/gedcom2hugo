package main

import (
	"os"

	"github.com/tektsu/gedcom2hugo/cmd"
	"github.com/urfave/cli"
)

const version = "0.0.0.1"

func main() {

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print only the version",
	}

	app := cli.NewApp()
	app.Name = "gedcom2hugo"
	app.Usage = "Generate Hugo input files from a GEDCOM file"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Enable verbose output",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debugging output",
		},
		cli.StringFlag{
			Name:  "gedcom, g",
			Usage: "Specify the input GEDCOM file",
		},
		cli.StringFlag{
			Name:  "project, p",
			Value: ".",
			Usage: "Specify the top level directory of the Hugo project",
		},
	}
	app.ArgsUsage = " "
	app.HideHelp = true

	app.Action = cmd.Generate

	app.Run(os.Args)
}
