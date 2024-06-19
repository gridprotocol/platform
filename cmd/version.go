package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const Version = "0.2.2"

var BuildFlag string

var VersionCmd = &cli.Command{
	Name:    "version",
	Usage:   "print platform version",
	Aliases: []string{"V"},
	Action: func(_ *cli.Context) error {
		fmt.Println(Version + "+" + BuildFlag)
		return nil
	},
}
