package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/satori/go.uuid"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	Apps  []App
	Hosts []Host
}

type App struct {
	Name  string
	Dir   string
	Hosts []string
}

type Host struct {
	Uid  string
	Ip   string
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

func LoadConfig() *Configuration {
	cfg := &Configuration{}
	cfg.readConfigurationFromDisk()
	return cfg
}

func (c *Configuration) addApp(a App) error {
	for _, ap := range c.Apps {
		if ap.Name == a.Name {
			return fmt.Errorf("App %s já existe", a.Name)
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
	return fmt.Errorf("App %s não foi encontrado", a)
}

func (c *Configuration) addHost(h Host) error {
	for _, ho := range c.Hosts {
		if ho.Ip == h.Ip && h.Port == ho.Port {
			return fmt.Errorf("Host %s já existe", h.Ip)
		}
	}
	h.Uid = uuid.NewV4().String()
	c.Hosts = append(c.Hosts, h)
	return nil
}

func (c *Configuration) findHostPosition(s string) (err error, i int) {
	if len(s) < 3 {
		err = fmt.Errorf("Informe ao menos 3 caracters do UID")
		return
	}
	for i = len(c.Hosts) - 1; i >= 0; i-- {
		h := c.Hosts[i]
		if strings.HasPrefix(h.Uid, s) {
			return
		}
	}
	err = fmt.Errorf("Host %s não encontrado", s)
	return
}

func (c *Configuration) findAppPosition(s string) (err error, i int) {
	for i = len(c.Apps) - 1; i >= 0; i-- {
		a := c.Apps[i]
		if strings.HasPrefix(a.Name, s) {
			return
		}
	}
	err = fmt.Errorf("App %s não encontrado", s)
	return
}

func (c *Configuration) removeHost(h string) error {
	if len(h) < 3 {
		return fmt.Errorf("Informe pelo menos 3 caracteres do UID")
	}
	for i := len(c.Hosts) - 1; i >= 0; i-- {
		host := c.Hosts[i]
		if host.Uid == h {
			c.Hosts = append(c.Hosts[:i], c.Hosts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Host %s não encontrado", h)
}

func (c *Configuration) addHostToApp(a string, h string) error {
	err, app := c.findAppPosition(a)
	if err != nil {
		return fmt.Errorf("App %s não encontrado", a)
	}
	err, host := c.findHostPosition(h)
	if err != nil {
		return fmt.Errorf("Host %s não encontrado", h)
	}
	for _, ho := range c.Apps[app].Hosts {
		if ho == c.Hosts[host].Uid {
			return fmt.Errorf("Host %s já esta associado ao app %s", c.Hosts[host].Ip, c.Apps[app].Name)
		}
	}
	c.Apps[app].Hosts = append(c.Apps[app].Hosts, c.Hosts[host].Uid)
	return nil
}

func (c *Configuration) removeHostFromApp(a string, h string) error {
	err, app := c.findAppPosition(a)
	if err != nil {
		return fmt.Errorf("App %s não encontrado", a)
	}
	err, host := c.findHostPosition(h)
	if err != nil {
		return fmt.Errorf("Host %s não encontrado", h)
	}
	for i := len(c.Apps[app].Hosts) - 1; i >= 0; i-- {
		ho := c.Apps[app].Hosts[i]
		if ho == c.Hosts[host].Uid {
			c.Apps[app].Hosts = append(c.Apps[app].Hosts[:i], c.Apps[app].Hosts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Host %s não foi encontrado no app %s", c.Hosts[host].Ip, c.Apps[app].Name)
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
			log.Fatalln("Ocorreu um erro ao ler o arquivo de configuração.")
		}
		file, _ = os.Create(path)
		configYaml, _ := yaml.Marshal(t)
		_, err = file.Write(configYaml)
		file, _ = os.Open(path)
	}
	byte, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Ocorreu um erro ao ler o arquivo de configuração.")
	}
	err = yaml.Unmarshal(byte, t)
	if err != nil {
		log.Fatalln("Formato invalido do arquivo de configuração", err)
	}
}

func (c *Configuration) saveConfigToDisk() {
	cfgFilePath := filepath.Join(cfgFolder, cfgFileName)
	c.saveYamlToDisk(c, cfgFilePath)
}

func (c *Configuration) saveYamlToDisk(i interface{}, path string) {
	h, err := yaml.Marshal(i)
	if err != nil {
		log.Fatalln("Não foi possivel fazer o parse das configurações", err)
	}
	ioutil.WriteFile(path, h, 0644)
}
