package main

import (
	"github.com/urfave/cli"
)

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:      "reload",
			Usage:     "Recria containers com novas configurações",
			ArgsUsage: "[app]",
			Before:    defineDockerHostCommand,
			Action:    reloadAction,
		},
		{
			Name:      "stop",
			Usage:     "Para os containers",
			ArgsUsage: "[app]",
			Before:    defineDockerHostCommand,
			Action:    stopAction,
		},
		{
			Name:      "deploy",
			Usage:     "Inicia o deploy da aplicação",
			ArgsUsage: "[app]",
			Before:    defineDockerHostCommand,
			Flags:     getDeployTeardownFlags(),
			Action:    deployAction,
		},
		{
			Name:      "teardown",
			Usage:     "Destroi o ambiente da aplicação",
			ArgsUsage: "[app]",
			Before:    defineDockerHostCommand,
			Flags:     getDeployTeardownFlags(),
			Action:    tearDownAction,
		},
		{
			Name:   "update-images",
			Usage:  "Faz update das imagens a partir do mirror",
			Before: defineDockerHostCommand,
			Action: forceUpdateImages,
		},
		{
			Name:            "docker",
			Usage:           "Executa comandos docker no node definido como Host",
			SkipFlagParsing: true,
			Before:          defineDockerHostCommand,
			Action:          executeDockerCmd,
		},
		{
			Name:            "compose",
			Usage:           "Executa comandos compose no node definido como Host",
			SkipFlagParsing: true,
			Before:          defineDockerHostCommand,
			Action:          executeDockerComposeCmd,
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
