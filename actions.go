package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

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
		dockerArgs = append(dockerArgs, "--tls", "-H", "tcp://"+host)
		dockerComposeArgs = append(dockerComposeArgs, "-H", "tcp://"+host)
	}
	if len(certsPath) > 0 {
		certsArgs := []string{"--tlscert=" + certsPath + "/cert.pem", "--tlskey=" + certsPath + "/key.pem"}
		dockerComposeArgs = append(dockerComposeArgs, certsArgs...)
		dockerArgs = append(dockerArgs, certsArgs...)
	}
	return nil
}

func deployAction(c *cli.Context) error {
	dockerCompose := parseDockerCompose()
	updateServicesImage(dockerCompose)
	args := []string{"up", "-d"}
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
			return fmt.Errorf("Falha ao obter a key %s do config", name)
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

var ListHostTemplate = `IP{{"\t\t\t"}}PORT{{range .Hosts}}
{{.Ip}}{{"\t\t\t"}}{{.Port}}{{end}}
`

var ListHostTeplateQuiet = `{{range .Hosts}}{{.Ip}}{{"\n"}}{{end}}`

func listHosts(c *cli.Context) error {
	cfg := LoadConfig()
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
	tp := ListHostTemplate
	if quiet {
		tp = ListHostTeplateQuiet
	}
	t := template.Must(template.New("ls").Parse(tp))
	err := t.Execute(w, cfg)
	if err != nil {
		return fmt.Errorf("Ocorreu um erro: %s", err)
	}
	w.Flush()
	return nil
}

func removeHost(c *cli.Context) error {
	args := c.Args()
	if len(args.First()) == 0 {
		return fmt.Errorf("\"brawl host rm\" requer um parametro")
	}
	cfg := LoadConfig()
	for _, a := range args {
		for _, h := range cfg.Hosts {
			if h.Ip == a {
				cfg.removeHost(h)
				fmt.Println(h.Ip)
			}
		}
	}
	cfg.saveConfigToDisk()
	return nil
}

func createHost(c *cli.Context) error {
	args := c.Args()
	host := Host{
		Ip:   args.Get(0),
		Port: args.Get(1),
	}
	if len(host.Ip) == 0 || len(host.Port) == 0 {
		return fmt.Errorf("\"brawl host create\" requer ip e porta como parametros")
	}
	cfg := LoadConfig()
	for _, h := range cfg.Hosts {
		if h.Ip == host.Ip {
			return fmt.Errorf("Host %s já está cadastrado", host.Ip)
		}
	}
	cfg.addHost(host)
	cfg.saveConfigToDisk()
	fmt.Println(host.Ip)
	return nil
}
