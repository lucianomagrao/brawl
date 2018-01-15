package actions

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func (act *Action) ExecuteDockerCmd(c *cli.Context) error {
	args := append(dockerArgs, c.Args()...)
	code := execDockerCommandAndWait(args...)
	if code > 0 {
		return act.CreateErrorMessage("Ocorreu um erro na execução, código de saida: %i", code)
	}
	return nil
}

func PushImage(img string) {
	ec := execDockerCommandAndWait("push", img)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel fazer o push da imagem %s", img)
	}
}

func PullImage(img string) {
	ec := execDockerCommandAndWait("pull", img)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel fazer o pull da imagem %s", img)
	}
}

func TagImage(img, tag string) {
	ec := execDockerCommandAndWait("tag", img, tag)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel criar a tag para a imagem %s", img)
	}
}

func RemoveImage(img string) {
	ec := execDockerCommandAndWait("rmi", img)
	if ec > 0 {
		log.Fatalf(color.RedString("Error: ")+"Não foi possivel remover a imagem %s", img)
	}
}

func execDockerCommandAndWait(args ...string) int {
	runArgs := append(dockerArgs, args...)
	return execCommandAndWait("docker", runArgs...)
}

func execCommandAndWait(location string, command ...string) int {
	cmd := exec.Command(location, command...)
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
