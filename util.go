package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/project"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

const (
	dockerRegistryLocalEnv  = "DOCKER_REGISTRY"
	dockerRegistryMirrorEnv = "DOCKER_REGISTRY_MIRROR"
)

var (
	dockerArgs, dockerComposeArgs, hosts []string
	dockerRegistry                       string = os.Getenv(dockerRegistryLocalEnv)
	dockerRegistryMirror                 string = os.Getenv(dockerRegistryMirrorEnv)
)

func createComposeProjectContext() project.Context {
	composeFile := getWorkingDir() + "/docker-compose.yml"
	return project.Context{
		ProjectName:  "brawl",
		ComposeFiles: []string{composeFile},
	}
}

func parseDockerCompose() *project.Project {
	composeFile := getWorkingDir() + "/docker-compose.yml"
	ctx := createComposeProjectContext()
	p := project.NewProject(&ctx, nil, &config.ParseOptions{})
	if err := p.Parse(); err != nil {
		log.Fatalf(color.RedString("Error: ")+"Falha ao executar o parse do arquivo %s: %v", composeFile, err)
	}
	if p.ServiceConfigs == nil {
		log.Fatalf(color.RedString("Error: ") + "Nenhum serviço encontrado, abortando")
	}
	return p
}

func getServiceVersionFromProperties(serviceName string) string {
	var envMap map[string]string
	envMap, err := godotenv.Read(getWorkingDir() + "/.env")
	if err != nil {
		log.Fatalf(color.RedString("Error: ")+"Falha ao extrair variaveis do envfile: %v", err)
	}
	version := envMap[strings.ToUpper(serviceName)+"_VERSION"]
	if len(version) == 0 {
		log.Fatalln(color.RedString("Error: ") + "Variavel VERSION não informada no .env")
	}
	return version
}

func updateServicesImage(dockerCompose *project.Project) {
	if len(dockerRegistry) == 0 {
		printInfoMessage("Variavel %s não foi setada, as imagens serão baixadas diretamente de %s podendo causar lentidão.\n", dockerRegistryLocalEnv, dockerRegistryMirror)
		dockerRegistry = dockerRegistryMirror
		os.Setenv(dockerRegistryLocalEnv, dockerRegistryMirror)
	}
	printInfoMessage("Fazendo cache das imagens docker no registry local --> \n")
	for _, name := range dockerCompose.ServiceConfigs.Keys() {
		createCacheServiceImage(name, dockerCompose)
	}
}

func createCacheServiceImage(name string, dockerCompose *project.Project) {
	service, ok := dockerCompose.ServiceConfigs.Get(name)
	if !ok {
		log.Fatalf(color.RedString("Error: ")+"Falha ao obter a key %s do config", name)
	}
	if len(service.Image) == 0 {
		printInfoMessage("Serviço %s não utiliza imagem. Continuando...", name)
		return
	}
	repl := strings.NewReplacer(
		"${", "",
		"}", "",
		strings.ToUpper(name)+"_VERSION", getServiceVersionFromProperties(name),
		dockerRegistryLocalEnv, dockerRegistry,
	)
	img := repl.Replace(service.Image)
	pl := execDockerCommandAndWait("pull", img)
	if pl > 0 && dockerRegistry != dockerRegistryMirror {
		mirr := strings.Replace(img, dockerRegistry, dockerRegistryMirror, -1)
		pullImage(mirr)
		tagImage(mirr, img)
		pushImage(img)
		removeImage(mirr)
	}
}

func pushImage(img string) {
	ec := execDockerCommandAndWait("push", img)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel fazer o push da imagem %s", img)
	}
}

func pullImage(img string) {
	ec := execDockerCommandAndWait("pull", img)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel fazer o pull da imagem %s", img)
	}
}

func tagImage(img, tag string) {
	ec := execDockerCommandAndWait("tag", img, tag)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel criar a tag para a imagem %s", img)
	}
}

func removeImage(img string) {
	ec := execDockerCommandAndWait("rmi", img)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel remover a imagem %s", img)
	}
}

func execDockerCommandAndWait(args ...string) int {
	runArgs := append(dockerArgs, args...)
	return execCommandAndWait("docker", runArgs...)
}

func execComposeCommandAndWait(args ...string) int {
	runArgs := append(dockerComposeArgs, args...)
	return execCommandAndWait("docker-compose", runArgs...)
}

func execCommandAndWait(location string, command ...string) int {
	cmd := exec.Command(location, command...)
	cmd.Dir = getWorkingDir()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	return 0
}

func printInfoMessage(m string, a ...interface{}) (int, error) {
	blue := color.New(color.FgCyan, color.Bold).SprintFunc()
	return fmt.Printf(blue("Info: ")+m, a...)
}

func throwErrorMessage(m string, a ...interface{}) error {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	return fmt.Errorf(red("Error: ") + fmt.Sprintf(m, a...))
}
