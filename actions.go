package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/urfave/cli"
	"regexp"
	"strings"
)

var ValidIpAddressRegex = `[0-9]+(?:\.[0-9]+){3}:[0-9]+`

func forceUpdateImages(c *cli.Context) error {
	dockerCompose := parseDockerCompose()
	updateServicesImage(dockerCompose)
	return nil
}

func executeDockerCmd(c *cli.Context) error {
	code := execDockerCommandAndWait(c.Args()...)
	if code > 0 {
		return fmt.Errorf("Ocorreu um erro na execução, código de saida: %i", code)
	}
	return nil
}

func executeDockerComposeCmd(c *cli.Context) error {
	code := execComposeCommandAndWait(c.Args()...)
	if code > 0 {
		return fmt.Errorf("Ocorreu um erro na execução, código de saida: %i", code)
	}
	return nil
}

func defineDockerHostCommand(c *cli.Context) error {
	cfg := LoadConfig()
	if len(app) > 0 {
		_, a := cfg.findAppPosition(app)
		for _, host := range cfg.Apps[a].Hosts {
			_, h := cfg.findHostPosition(host)
			ho := cfg.Hosts[h]
			hosts = append(hosts, ho.Ip+":"+ho.Port)
		}
		workingDir = cfg.Apps[a].Dir
	}
	if len(host) > 0 {
		rp := regexp.MustCompile(ValidIpAddressRegex)
		tcp := "tcp://"
		if rp.MatchString(host) {
			tcp += host
		} else {
			err, h := cfg.findHostPosition(host)
			if err != nil {
				return err
			}
			ho := cfg.Hosts[h]
			tcp += ho.Ip + ":" + ho.Port
		}
		dockerArgs = append(dockerArgs, "-H", tcp)
		dockerComposeArgs = append(dockerComposeArgs, "-H", tcp)
	}
	if insecure {
		os.Unsetenv("DOCKER_CERT_PATH")
	} else {
		dockerArgs = append(dockerArgs, "--tls")
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
	for _, h := range hosts {
		args := []string{"-H", h, "up", "-d"}
		args = append(args, c.Args()...)
		execComposeCommandAndWait(args...)
	}
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

var ListHostTemplate = `UID{{"\t\t\t"}}IP{{"\t\t\t"}}PORT{{range .Hosts}}
{{.Uid}}{{"\t\t\t"}}{{.Ip}}{{"\t\t\t"}}{{.Port}}{{end}}
`

var ListHostTeplateQuiet = `{{range .Hosts}}{{.Uid}}{{"\n"}}{{end}}`

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
		err := cfg.removeApp(a)
		if err != nil {
			return err
		}
		fmt.Println(a)
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
		err := cfg.removeHost(a)
		if err != nil {
			return err
		}
		fmt.Println(a)
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
	fmt.Println(host.Uid)
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
		return fmt.Errorf("\"brawl app manage add-host\" requer 2 parametros.")
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
		return fmt.Errorf("\"brawl app manage rem-host\" requer 2 parametros.")
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
