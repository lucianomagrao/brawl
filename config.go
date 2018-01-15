package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Services map[string]Node `yaml:",inline"`
}

type Node struct {
	Nodes map[string]NodeConfig `yaml:",inline"`
}

type NodeConfig struct {
	IP      string `yaml:"ip"`
	Version string `yaml:"version"`
}

const (
	nodesFileName = "nodes.yml"
)

//LoadConfig Loads configuration from config file on disk
func LoadConfig() *Service {
	cfg := &Service{}
	cfg.readConfigFileFromDisk(nodesFileName)
	return cfg
}

func (c *Service) readConfigFileFromDisk(fileName string) {
	path := filepath.Join(getWorkingDir(), fileName)
	byte, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(color.RedString("Error: ") + "Ocorreu um erro ao ler o arquivo de configuração.")
	}
	err = yaml.Unmarshal(byte, &c)
	if err != nil {
		log.Fatalln(color.RedString("Error: ")+"Formato invalido do arquivo de configuração", err)
	}
}
