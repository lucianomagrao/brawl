package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/urfave/cli"
	"strings"
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

var ListAppTemplate = `NAME{{"\t\t\t"}}DIR{{"\t\t\t"}}HOSTS{{range .Apps}}
{{.Name}}{{"\t\t\t"}}{{.Dir}}{{"\t\t\t"}}{{join .Hosts ", "}}{{end}}
`

var ListAppTeplateQuiet = `{{range .Apps}}{{.Name}}{{"\n"}}{{end}}`

func listHosts(c *cli.Context) error {
	return list(ListHostTemplate, ListHostTeplateQuiet)
}

func listApps(c *cli.Context) error {
	return list(ListAppTemplate, ListAppTeplateQuiet)
}

func list(tmpl string, qtmpl string) error {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}
	cfg := LoadConfig()
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
	tp := tmpl
	if quiet {
		tp = qtmpl
	}
	t := template.Must(template.New("ls").Funcs(funcMap).Parse(tp))
	err := t.Execute(w, cfg)
	if err != nil {
		return fmt.Errorf("Ocorreu um erro: %s", err)
	}
	w.Flush()
	return nil
}

func removeApp(c *cli.Context) error {
	args := c.Args()
	if len(args.First()) == 0 {
		return fmt.Errorf("\"brawl app rm\" requer um parametro")
	}
	cfg := LoadConfig()
	for _, a := range args {
		for _, h := range cfg.Apps {
			if h.Name == a {
				err := cfg.removeApp(h)
				if err != nil {
					return err
				}
				fmt.Println(h.Name)
			}
		}
	}
	cfg.saveConfigToDisk()
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
				err := cfg.removeHost(h)
				if err != nil {
					return err
				}
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
	err := cfg.addHost(host)
	if err != nil {
		return err
	}
	cfg.saveConfigToDisk()
	fmt.Println(host.Ip)
	return nil
}

func createApp(c *cli.Context) error {
	args := c.Args()
	app := App{
		Name: args.Get(0),
		Dir:  args.Get(1),
	}
	if len(app.Name) == 0 || len(app.Dir) == 0 {
		return fmt.Errorf("\"brawl app create\" requer nome e diretório como parametros")
	}
	cfg := LoadConfig()
	err := cfg.addApp(app)
	if err != nil {
		return err
	}
	cfg.saveConfigToDisk()
	fmt.Println(app.Name)
	return nil
}

func addHostToApp(c *cli.Context) error {
	args := c.Args()
	a := args.Get(0)
	h := args.Get(1)
	if len(a) == 0 || len(h) == 0 {
		return fmt.Errorf("\"brawl app add-host\" requer 2 parametros.")
	}
	cfg := LoadConfig()
	err := cfg.addHostToApp(a, h)
	if err != nil {
		return err
	}
	fmt.Println(h)
	cfg.saveConfigToDisk()
	return nil
}

func removeHostFromApp(c *cli.Context) error {
	args := c.Args()
	a := args.Get(0)
	h := args.Get(1)
	if len(a) == 0 || len(h) == 0 {
		return fmt.Errorf("\"brawl app rem-host\" requer 2 parametros.")
	}
	cfg := LoadConfig()
	err := cfg.removeHostFromApp(a, h)
	if err != nil {
		return err
	}
	fmt.Println(h)
	cfg.saveConfigToDisk()
	return nil
}
