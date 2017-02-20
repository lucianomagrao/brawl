package main

import (
	"os"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "brawl"
	app.Version = "1.0.2"
	app.Usage = "Deploy de aplicações docker."

	app.Authors = []cli.Author{
		{
			Name:  "Luciano Mores",
			Email: "luciano.mores@softplan.com.br",
		},
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Exibe a versão do sistema",
	}

	app.Flags = getFlags()
	app.Commands = getCommands()

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}
