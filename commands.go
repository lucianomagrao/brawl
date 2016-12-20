package main

import (
	"github.com/urfave/cli"
)

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:      "deploy",
			Usage:     "Inicia o deploy da aplicação",
			ArgsUsage: "[aplicações...]",
			Before:    defineDockerHostCommand,
			Action:    deployAction,
		},
		{
			Name:   "update-images",
			Usage:  "Faz update das imagens a partir do mirror",
			Before: defineDockerHostCommand,
			Action: forceUpdateImages,
		},
		{
			Name:            "docker",
			Usage:           "Executa comandos docker na maquina definida como Host",
			SkipFlagParsing: true,
			Before:          defineDockerHostCommand,
			Action:          executeDockerCmd,
		},
		{
			Name:            "compose",
			Usage:           "Executa comandos compose na maquina definida como Host",
			SkipFlagParsing: true,
			Before:          defineDockerHostCommand,
			Action:          executeDockerComposeCmd,
		},
		{
			Name:  "host",
			Usage: "Gerencia hosts",
			Subcommands: []cli.Command{
				{
					Name:   "ls",
					Usage:  "Lista hosts cadastrados",
					Action: listHosts,
					Flags:  getLsFlags(),
				},
				{
					Name:      "rm",
					ArgsUsage: "[hosts...]",
					Usage:     "Remove host",
					Action:    removeHost,
				},
				{
					Name:      "create",
					ArgsUsage: "[ip] [porta]",
					Usage:     "Cria novo host",
					Action:    createHost,
				},
			},
		},
		{
			Name:  "app",
			Usage: "Gerencia apps",
			Subcommands: []cli.Command{
				{
					Name:   "ls",
					Usage:  "Lista apps cadastrados",
					Action: listApps,
					Flags:  getLsFlags(),
				},
				{
					Name:      "rm",
					ArgsUsage: "[apps...]",
					Usage:     "Remove apps",
					Action:    removeApp,
				},
				{
					Name:      "create",
					ArgsUsage: "[name] [dir]",
					Usage:     "Cria novo app",
					Action:    createApp,
				},
				{
					Name:      "manage",
					ArgsUsage: "[app]",
					Usage:     "Gerencia app",
					Subcommands: []cli.Command{
						{
							Name:      "add-host",
							ArgsUsage: "[host]",
							Usage:     "Adiciona host ao app",
							Action:    addHostToApp,
						},
						{
							Name:      "rm-host",
							ArgsUsage: "[host]",
							Usage:     "Remove host do app",
							Action:    removeHostFromApp,
						},
					},
				},
			},
		},
		{
			Name:   "versions",
			Usage:  "Exibe as versões das aplicações",
			Before: defineDockerHostCommand,
			Action: showVersionsAction,
		},
		{
			Name:      "help",
			Aliases:   []string{"h"},
			Usage:     "Exibe a lista de commandos ou a ajuda de um comando",
			ArgsUsage: "[commando]",
			Action:    showHelpAction,
		},
	}
}
