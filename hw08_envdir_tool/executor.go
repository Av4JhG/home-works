package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

var ErrCommandInfoNotFound = errors.New("no command provided")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return int(syscall.EINVAL)
	}

	setEnvVars(env)

	command := getCommand(cmd)
	command.Run()

	return command.ProcessState.ExitCode()
}

func setEnvVars(env Environment) {
	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
			continue
		}

		os.Setenv(name, value.Value)
	}
}

func getCommand(cmd []string) *exec.Cmd {
	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command
}
