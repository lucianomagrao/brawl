package main

import (
	"github.com/urfave/cli"
	"os"
	"sort"
)

func main() {
	app := cli.NewApp()
	app.Name = "Brawl!"
	app.Version = "1.0.0"
	app.Usage = "Deploy de aplicações docker da Softplan."

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

	cfg := LoadConfig()
	app.Run(os.Args)
	cfg.saveConfigToDisk()
}
