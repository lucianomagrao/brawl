package actions

import (
	"github.com/docker/libcompose/project"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func (act *Action) UpdateServicesImage() {
	dockerCompose := act.parseDockerCompose()
	if len(dockerRegistry) == 0 {
		act.PrintInfoMessage("Variavel %s não foi setada, as imagens serão baixadas diretamente de %s podendo causar lentidão.\n", dockerRegistryLocalEnv, dockerRegistryMirror)
		dockerRegistry = dockerRegistryMirror
		os.Setenv(dockerRegistryLocalEnv, dockerRegistryMirror)
	}
	act.PrintInfoMessage("Fazendo cache das imagens docker no registry local --> \n")
	for _, name := range dockerCompose.ServiceConfigs.Keys() {
		act.createCacheServiceImage(name, dockerCompose)
	}
}

func (act *Action) createCacheServiceImage(name string, dockerCompose *project.Project) {
	service, ok := dockerCompose.ServiceConfigs.Get(name)
	if !ok {
		log.Fatalf(color.RedString("Error: ")+"Falha ao obter a key %s do config", name)
	}
	if len(service.Image) == 0 {
		act.PrintInfoMessage("Serviço %s não utiliza imagem. Continuando...", name)
		return
	}
	repl := strings.NewReplacer(
		"${", "",
		"}", "",
		strings.ToUpper(name)+"_VERSION", act.getServiceVersionFromProperties(name),
		dockerRegistryLocalEnv, dockerRegistry,
	)
	img := repl.Replace(service.Image)
	pl := execDockerCommandAndWait("pull", img)
	if pl > 0 && dockerRegistry != dockerRegistryMirror {
		mirr := strings.Replace(img, dockerRegistry, dockerRegistryMirror, -1)
		PullImage(mirr)
		TagImage(mirr, img)
		PushImage(img)
		RemoveImage(mirr)
	}
}

func (act *Action) getServiceVersionFromProperties(serviceName string) string {
	var envMap map[string]string
	envMap, err := godotenv.Read(act.workingDir + "/" + serviceName + ".properties")
	if err != nil {
		log.Fatalf(color.RedString("Error: ")+"Falha ao extrair variaveis do envfile: %v", err)
	}
	version := envMap[strings.ToUpper(serviceName)+"_VERSION"]
	if len(version) == 0 {
		log.Fatalln(color.RedString("Error: ") + "Variavel VERSION não informada no .env")
	}
	return version
}
