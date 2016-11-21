package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

func forceUpdateImages(c *cli.Context) error {
	dockerCompose := parseDockerCompose()
	updateServicesImage(dockerCompose)
	return nil
}

func executeDockerCmd(c *cli.Context) error {
	execDockerCommandAndWait(c.Args()...)
	return nil
}

func executeDockerComposeCmd(c *cli.Context) error {
	execComposeCommandAndWait(c.Args()...)
	return nil
}

func defineDockerHostCommand(c *cli.Context) error {
	if len(host) > 0 {
		dockerArgs = append(
			dockerArgs,
			"--tls",
			"-H",
			"tcp://"+host,
		)
		dockerComposeArgs = append(
			dockerComposeArgs,
			"-H",
			"tcp://"+host,
		)
	}
	if len(certsPath) > 0 {
		certsArgs := []string{
			"--tlscert=" + certsPath + "/cert.pem",
			"--tlskey=" + certsPath + "/key.pem",
		}
		dockerComposeArgs = append(dockerComposeArgs, certsArgs...)
		dockerArgs = append(dockerArgs, certsArgs...)
	}
	return nil
}

func deployAction(c *cli.Context) error {
	dockerCompose := parseDockerCompose()
	updateServicesImage(dockerCompose)
	args := []string{
		"up",
		"-d",
	}
	args = append(args, c.Args()...)
	execComposeCommandAndWait(args...)
	return nil
}

func showVersionsAction(c *cli.Context) error {
	dockerCompose := parseDockerCompose()
	fmt.Printf("\nVersões das aplicações\n\n")
	for _, name := range dockerCompose.ServiceConfigs.Keys() {
		service, ok := dockerCompose.ServiceConfigs.Get(name)
		if !ok {
			log.Fatalf("Falha ao obter a key %s do config", name)
		}
		if len(service.Image) == 0 {
			continue
		}
		fmt.Printf("%s=%s\n", name, getServiceVersionFromProperties(name))
	}
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
