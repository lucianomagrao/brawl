package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"regexp"
	"strings"

	"github.com/urfave/cli"
)

//ValidIPAddressRegex Regex to validate IP addresses
var ValidIPAddressRegex = `[0-9]+(?:\.[0-9]+){3}:[0-9]+`

//ListHostTemplate Template for Host listing
var ListHostTemplate = `UID{{"\t\t\t"}}IP{{"\t\t\t"}}PORT{{range .Hosts}}
{{.UID}}{{"\t\t\t"}}{{.IP}}{{"\t\t\t"}}{{.Port}}{{end}}
`

//ListHostTeplateQuiet Template for quiet Host listing
var ListHostTeplateQuiet = `{{range .Hosts}}{{.UID}}{{"\n"}}{{end}}`

//ListAppTemplate Template for App listing
var ListAppTemplate = `NAME{{"\t\t\t"}}DIR{{"\t\t\t"}}HOSTS{{range .Apps}}
{{.Name}}{{"\t\t\t"}}{{.Dir}}{{"\t\t\t"}}{{join .Hosts ", "}}{{end}}
`

//ListAppTeplateQuiet Template for quiet App listing
var ListAppTeplateQuiet = `{{range .Apps}}{{.Name}}{{"\n"}}{{end}}`

func forceUpdateImages(c *cli.Context) error {
	dockerArgs = append(dockerArgs, "-H", hosts[0])
	dockerComposeArgs = append(dockerComposeArgs, "-H", hosts[0])
	dockerCompose := parseDockerCompose()
	updateServicesImage(dockerCompose)
	return nil
}

func executeDockerCmd(c *cli.Context) error {
	for _, host := range hosts {
		args := append(dockerArgs, "-H", host)
		args = append(args, c.Args()...)
		printInfoMessage("Executando comando no host --> %s \n", host)
		code := execDockerCommandAndWait(args...)
		if code > 0 {
			return createErrorMessage("Ocorreu um erro na execução, código de saida: %i", code)
		}
	}
	return nil
}

func executeDockerComposeCmd(c *cli.Context) error {
	for _, host := range hosts {
		args := append(dockerComposeArgs, "-H", host)
		args = append(args, c.Args()...)
		printInfoMessage("Executando comando no host --> %s \n", host)
		code := execDockerCommandAndWait(args...)
		if code > 0 {
			return createErrorMessage("Ocorreu um erro na execução, código de saida: %i", code)
		}
	}
	return nil
}

func defineDockerHostCommand(c *cli.Context) error {
	cfg := LoadConfig()
	var ho Host
	if len(app) > 0 {
		a, _ := cfg.findAppPosition(app)
		for _, host := range cfg.Apps[a].Hosts {
			h, _ := cfg.findHostPosition(host)
			ho = cfg.Hosts[h]
			hosts = append(hosts, ho.IP+":"+ho.Port)
		}
		workingDir = cfg.Apps[a].Dir
	}
	if len(host) > 0 {
		rp := regexp.MustCompile(ValidIPAddressRegex)
		var addr string
		if rp.MatchString(host) {
			addr = host
		} else {
			h, err := cfg.findHostPosition(host)
			if err != nil {
				return err
			}
			addr = cfg.Hosts[h].IP + ":" + cfg.Hosts[h].Port
		}
		if len(app) > 0 {
			hosts = []string{addr}
		}
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
		printInfoMessage("Iniciando deploy no host --> %s \n", h)
		execComposeCommandAndWait(args...)
	}
	return nil
}

func stopAction(c *cli.Context) error {
	for _, h := range hosts {
		args := []string{"-H", h, "down"}
		args = append(args, c.Args()...)
		printInfoMessage("Parando containers do host --> %s \n", h)
		execComposeCommandAndWait(args...)
	}
	return nil
}

func reloadAction(c *cli.Context) error {
	stopAction(c)
	for _, h := range hosts {
		args := []string{"-H", h, "up", "-d"}
		args = append(args, c.Args()...)
		printInfoMessage("Inciando containers do host --> %s \n", h)
		execComposeCommandAndWait(args...)
	}
	return nil
}

func showVersionsAction(c *cli.Context) error {
	dockerCompose := parseDockerCompose()
	printInfoMessage("\nVersões das aplicações\n\n")
	for _, name := range dockerCompose.ServiceConfigs.Keys() {
		service, ok := dockerCompose.ServiceConfigs.Get(name)
		if !ok {
			return createErrorMessage("Falha ao obter a key %s do config", name)
		}
		if len(service.Image) == 0 {
			continue
		}
		printInfoMessage("%s=%s\n", name, getServiceVersionFromProperties(name))
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
		return createErrorMessage("Ocorreu um erro: %s", err)
	}
	w.Flush()
	return nil
}

func removeApp(c *cli.Context) error {
	args := c.Args()
	if len(args.First()) == 0 {
		return cli.ShowSubcommandHelp(c)
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
		return cli.ShowSubcommandHelp(c)
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
		IP:   args.Get(0),
		Port: args.Get(1),
	}
	if len(host.IP) == 0 || len(host.Port) == 0 {
		return cli.ShowSubcommandHelp(c)
	}
	cfg := LoadConfig()
	h, err := cfg.addHost(host)
	if err != nil {
		return err
	}
	cfg.saveConfigToDisk()
	fmt.Println(h.UID)
	return nil
}

func createApp(c *cli.Context) error {
	args := c.Args()
	app := App{
		Name: args.Get(0),
		Dir:  args.Get(1),
	}
	if len(app.Name) == 0 || len(app.Dir) == 0 {
		return cli.ShowSubcommandHelp(c)
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
		return cli.ShowSubcommandHelp(c)
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
		return cli.ShowSubcommandHelp(c)
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
