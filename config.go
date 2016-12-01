package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	Ip   string
	Port string
}

const (
	cfgFolderName = ".brawl"
	cfgFileName   = "config.json"
)

var (
	cfgFolder string = filepath.Join(os.Getenv("HOME"), cfgFolderName)
)

func LoadConfig() *Configuration {
	cfg := &Configuration{}
	cfg.readConfigurationFromDisk()
	return cfg
}

func (c *Configuration) addApp(a App) {
	c.Apps = append(c.Apps, a)
}

func (c *Configuration) removeApp(a App) {
	for i := len(c.Apps) - 1; i >= 0; i-- {
		ap := c.Apps[i]
		if ap.Name == a.Name {
			c.Apps = append(c.Apps[:i], c.Apps[i+1:]...)
			return
		}
	}
}

func (c *Configuration) addHost(h Host) {
	c.Hosts = append(c.Hosts, h)
}

func (c *Configuration) removeHost(h Host) {
	for i := len(c.Hosts) - 1; i >= 0; i-- {
		host := c.Hosts[i]
		if host.Ip == h.Ip {
			c.Hosts = append(c.Hosts[:i], c.Hosts[i+1:]...)
			return
		}
	}
}

func (a *App) addHost(h Host) {
	a.Hosts = append(a.Hosts, h.Ip)
}

func (a *App) removeHost(h Host) {
	for i := len(a.Hosts) - 1; i >= 0; i-- {
		host := a.Hosts[i]
		if host == h.Ip {
			a.Hosts = append(a.Hosts[:i], a.Hosts[i+1:]...)
			return
		}
	}
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
		configJson, _ := json.MarshalIndent(t, "", " ")
		_, err = file.Write(configJson)
		file, _ = os.Open(path)
	}
	cfgDeco := json.NewDecoder(file)
	ec := cfgDeco.Decode(t)
	if ec != nil {
		log.Fatalln("Formato invalido do arquivo de configuração", ec)
	}
}

func (c *Configuration) saveConfigToDisk() {
	cfgFilePath := filepath.Join(cfgFolder, cfgFileName)
	c.saveJsonToDisk(c, cfgFilePath)
}

func (c *Configuration) saveJsonToDisk(i interface{}, path string) {
	h, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		log.Fatalln("Não foi possivel fazer o parse das configurações", err)
	}
	ioutil.WriteFile(path, h, 0644)
}
