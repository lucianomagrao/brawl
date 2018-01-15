package actions

import (
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/project"
	"github.com/fatih/color"
	"log"
)

func (act *Action) ExecComposeCommandAndWait(args ...string) int {
	runArgs := append(dockerComposeArgs, args...)
	return execCommandAndWait("docker-compose", runArgs...)
}

func (act *Action) createComposeProjectContext() project.Context {
	composeFile := act.workingDir + "/docker-compose.yml"
	return project.Context{
		ProjectName:  "brawl",
		ComposeFiles: []string{composeFile},
	}
}

func (act *Action) parseDockerCompose() *project.Project {
	ctx := act.createComposeProjectContext()
	p := project.NewProject(&ctx, nil, &config.ParseOptions{})
	if err := p.Parse(); err != nil {
		log.Fatalf(color.RedString("Error: ")+"Falha ao executar o parse do arquivo docker-compose.yml: %v", err)
	}
	if p.ServiceConfigs == nil {
		log.Fatalf(color.RedString("Error: ") + "Nenhum serviço encontrado, abortando")
	}
	return p
}
