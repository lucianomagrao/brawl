package main

import (
	"regexp"

	"fmt"
	"github.com/leslau/brawl/actions"
	"github.com/urfave/cli"
)

var (
	hosts               []string
	ValidIPAddressRegex = `[0-9]+(?:\.[0-9]+){3}:[0-9]+`
	action              actions.Action
)

func forceUpdateImages(c *cli.Context) error {
	action, err := actions.NewAction(hosts, getWorkingDir(), insecure, getCertsPath())
	if err != nil {
		panic(err)
	}
	action.UpdateServicesImage()
	return nil
}

func executeDockerCmd(c *cli.Context) error {
	code := action.ExecComposeCommandAndWait(c.Args()...)
	if code > 0 {
		return action.CreateErrorMessage("Ocorreu um erro na execução, código de saida: %i", code)
	}
	return nil
}

func executeDockerComposeCmd(c *cli.Context) error {
	code := action.ExecComposeCommandAndWait(c.Args()...)
	if code > 0 {
		return action.CreateErrorMessage("Ocorreu um erro na execução, código de saida: %i", code)
	}
	return nil
}

func defineDockerHostCommand(c *cli.Context) error {
	cfg := LoadConfig()
	fmt.Println(cfg.Services)
	for service, nodes := range cfg.Services {
		fmt.Printf("Service: [%s] \n", service)
		for node, config := range nodes.Nodes {
			fmt.Printf("Node: [%s] \n", node)
			fmt.Printf("Ip: [%s] \n", config.IP)
			if len(config.Version) > 0 {
				fmt.Printf("Version: [%s] \n", config.Version)
			}
		}
	}
	if len(host) > 0 {
		rp := regexp.MustCompile(ValidIPAddressRegex)
		var addr string
		if rp.MatchString(host) {
			addr = host
		}
		hosts = []string{addr}
	} else {

	}
	return nil
}

func deployAction(c *cli.Context) error {
	action, err := actions.NewAction(hosts, getWorkingDir(), insecure, getCertsPath())
	if err != nil {
		panic(err)
	}
	action.DeployService()
	return nil
}

func tearDownAction(c *cli.Context) error {
	action, err := actions.NewAction(hosts, getWorkingDir(), insecure, getCertsPath())
	if err != nil {
		panic(err)
	}
	action.TearDownService()
	return nil
}

func stopAction(c *cli.Context) error {
	action.ExecComposeCommandAndWait("stop")
	return nil
}

func reloadAction(c *cli.Context) error {
	action.ExecComposeCommandAndWait("restart")
	return nil
}

func showVersionsAction(c *cli.Context) error {
	return nil
}

func showHelpAction(c *cli.Context) error {
	args := c.Args()
	if args.Present() {
		return cli.ShowCommandHelp(c, args.First())
	}
	cli.ShowAppHelp(c)
	return nil
}
