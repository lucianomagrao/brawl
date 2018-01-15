package actions

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	dockerRegistryLocalEnv  = "DOCKER_REGISTRY"
	dockerRegistryMirrorEnv = "DOCKER_REGISTRY_MIRROR"
)

var (
	dockerArgs, dockerComposeArgs []string
	dockerRegistryMirror          string = os.Getenv(dockerRegistryMirrorEnv)
	dockerRegistry                string = os.Getenv(dockerRegistryLocalEnv)
)

type Action struct {
	// Hosts to execute the action
	hosts []string
	// Dir to refer docker-compose file
	workingDir string
	// Define if communication with docker host is secure or not
	insecure  bool
	certsPath string
}

func NewAction(hosts []string, workingDir string, insecure bool, certsPath string) (*Action, error) {
	if !insecure {
		dockerArgs = append(dockerArgs, "--tls")
	}
	if len(certsPath) > 0 {
		certsArgs := []string{"--tlscert=" + certsPath + "/cert.pem", "--tlskey=" + certsPath + "/key.pem"}
		dockerComposeArgs = append(dockerComposeArgs, certsArgs...)
		dockerArgs = append(dockerArgs, certsArgs...)
	}
	return &Action{
		hosts:      hosts,
		workingDir: workingDir,
		insecure:   insecure,
		certsPath:  certsPath,
	}, nil
}

func (act *Action) PrintInfoMessage(m string, a ...interface{}) (int, error) {
	blue := color.New(color.FgCyan, color.Bold).SprintFunc()
	return fmt.Printf(blue("Info: ")+m, a...)
}

func (act *Action) CreateErrorMessage(m string, a ...interface{}) error {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	return fmt.Errorf(red("Error: ") + fmt.Sprintf(m, a...))
}
