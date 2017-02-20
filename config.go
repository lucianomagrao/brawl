package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/satori/go.uuid"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

//Configuration config file struct
type Configuration struct {
	Apps  []App
	Hosts []Host
}

//App App struct
type App struct {
	Name  string
	Dir   string
	Hosts []string
}

//Host Host struct
type Host struct {
	UID  string
	IP   string
	Port string
}

const (
	cfgFolderName = ".brawl"
	cfgFileName   = "config.yaml"
)

var (
	home, _          = homedir.Dir()
	cfgFolder string = filepath.Join(home, cfgFolderName)
)

//LoadConfig Loads configuration from config file on disk
func LoadConfig() *Configuration {
	cfg := &Configuration{}
	cfg.readConfigurationFromDisk()
	return cfg
}

func (c *Configuration) addApp(a App) error {
	for _, ap := range c.Apps {
		if ap.Name == a.Name {
			return createErrorMessage("App %s já existe", a.Name)
		}
	}
	c.Apps = append(c.Apps, a)
	return nil
}

func (c *Configuration) removeApp(a string) error {
	for i := len(c.Apps) - 1; i >= 0; i-- {
		ap := c.Apps[i]
		if ap.Name == a {
			c.Apps = append(c.Apps[:i], c.Apps[i+1:]...)
			return nil
		}
	}
	return createErrorMessage("App %s não foi encontrado", a)
}

func (c *Configuration) addHost(h Host) (Host, error) {
	for _, ho := range c.Hosts {
		if ho.IP == h.IP && h.Port == ho.Port {
			return h, createErrorMessage("Host %s já existe", h.IP)
		}
	}
	h.UID = uuid.NewV4().String()
	c.Hosts = append(c.Hosts, h)
	return h, nil
}

func (c *Configuration) findHostPosition(s string) (i int, err error) {
	if len(s) < 3 {
		err = createErrorMessage("Informe ao menos 3 caracters do UID")
		return
	}
	for i = len(c.Hosts) - 1; i >= 0; i-- {
		h := c.Hosts[i]
		if strings.HasPrefix(h.UID, s) {
			return
		}
	}
	err = createErrorMessage("Host %s não encontrado", s)
	return
}

func (c *Configuration) findAppPosition(s string) (i int, err error) {
	for i = len(c.Apps) - 1; i >= 0; i-- {
		a := c.Apps[i]
		if strings.HasPrefix(a.Name, s) {
			return
		}
	}
	err = createErrorMessage("App %s não encontrado", s)
	return
}

func (c *Configuration) removeHost(h string) error {
	if len(h) < 3 {
		return createErrorMessage("Informe pelo menos 3 caracteres do UID")
	}
	for i := len(c.Hosts) - 1; i >= 0; i-- {
		host := c.Hosts[i]
		if strings.HasPrefix(host.UID, h) {
			c.Hosts = append(c.Hosts[:i], c.Hosts[i+1:]...)
			return nil
		}
	}
	return createErrorMessage("Host %s não encontrado", h)
}

func (c *Configuration) addHostToApp(a string, h string) error {
	app, err := c.findAppPosition(a)
	if err != nil {
		return createErrorMessage("App %s não encontrado", a)
	}
	host, err := c.findHostPosition(h)
	if err != nil {
		return createErrorMessage("Host %s não encontrado", h)
	}
	for _, ho := range c.Apps[app].Hosts {
		if ho == c.Hosts[host].UID {
			return createErrorMessage("Host %s já esta associado ao app %s", c.Hosts[host].IP, c.Apps[app].Name)
		}
	}
	c.Apps[app].Hosts = append(c.Apps[app].Hosts, c.Hosts[host].UID)
	return nil
}

func (c *Configuration) removeHostFromApp(a string, h string) error {
	app, err := c.findAppPosition(a)
	if err != nil {
		return createErrorMessage("App %s não encontrado", a)
	}
	host, err := c.findHostPosition(h)
	if err != nil {
		return createErrorMessage("Host %s não encontrado", h)
	}
	for i := len(c.Apps[app].Hosts) - 1; i >= 0; i-- {
		ho := c.Apps[app].Hosts[i]
		if ho == c.Hosts[host].UID {
			c.Apps[app].Hosts = append(c.Apps[app].Hosts[:i], c.Apps[app].Hosts[i+1:]...)
			return nil
		}
	}
	return createErrorMessage("Host %s não foi encontrado no app %s", c.Hosts[host].IP, c.Apps[app].Name)
}

func (c *Configuration) readConfigurationFromDisk() {
	os.MkdirAll(cfgFolder, os.ModeAppend)
	c.readConfigFileFromDisk(cfgFileName, c)
}

func (c *Configuration) readConfigFileFromDisk(fileName string, t interface{}) {
	path := filepath.Join(cfgFolder, fileName)
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln(color.RedString("Error: ") + "Ocorreu um erro ao ler o arquivo de configuração.")
		}
		file, _ = os.Create(path)
		configYaml, _ := yaml.Marshal(t)
		_, err = file.Write(configYaml)
		if err != nil {
			log.Fatalln(color.RedString("Error: ") + "Ocorreu um erro ao criar o arquivo de configuração.")
		}
		file, _ = os.Open(path)
	}
	byte, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(color.RedString("Error: ") + "Ocorreu um erro ao ler o arquivo de configuração.")
	}
	err = yaml.Unmarshal(byte, t)
	if err != nil {
		log.Fatalln(color.RedString("Error: ")+"Formato invalido do arquivo de configuração", err)
	}
}

func (c *Configuration) saveConfigToDisk() {
	cfgFilePath := filepath.Join(cfgFolder, cfgFileName)
	c.saveYamlToDisk(c, cfgFilePath)
}

func (c *Configuration) saveYamlToDisk(i interface{}, path string) {
	h, err := yaml.Marshal(i)
	if err != nil {
		log.Fatalln(color.RedString("Error: ")+"Não foi possivel fazer o parse das configurações", err)
	}
	ioutil.WriteFile(path, h, 0644)
}
