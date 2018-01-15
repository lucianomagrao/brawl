package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	workingDir, host, certsPath string
	insecure, force, all        bool
)

func getFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "folder, f",
			Value:       getDefaultWorkingDir(),
			Usage:       "Define a pasta de configuração a ser utilizada, caso não informado utiliza a pasta atual",
			Destination: &workingDir,
		},
		cli.StringFlag{
			Name:        "host, H",
			Usage:       "Define o host a ser utilizado",
			Destination: &host,
		},
		cli.StringFlag{
			Name:        "certs, c",
			Usage:       "Caminho onde contem os arquivos de certificado e key para acesso ao Host definido",
			Destination: &certsPath,
		},
		cli.BoolFlag{
			Name:        "insecure, i",
			Usage:       "Define se é insecure a comunicação com o host",
			Destination: &insecure,
		},
	}
}

func getDeployTeardownFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Ignora confirmação da ação de deploy",
			Destination: &force,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "Define se deve ser feito deploy em todos nodes",
			Destination: &all,
		},
	}
}

func getDefaultWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ocorreu um erro ao recuperar o a pasta local de execução %v", err)
	}
	return pwd
}

func getWorkingDir() string {
	_, err := os.Stat(workingDir)
	if os.IsNotExist(err) {
		log.Fatalf("Diretório informado não existe: %v", workingDir)
	}
	return workingDir
}

func getCertsPath() string {
	_, err := os.Stat(certsPath)
	if os.IsNotExist(err) {
		log.Fatalf("Diretório informado não existe: %v", certsPath)
	}
	return certsPath
}
