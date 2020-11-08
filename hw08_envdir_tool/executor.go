package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := cmd[0]
	c := exec.Command(command, cmd[1:]...)
	c.Env = os.Environ()
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr

	for k, v := range env {
		newEnv := k + "=" + v
		c.Env = append(c.Env, newEnv)
	}

	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		ok := errors.As(err, &exitErr)
		if ok {
			return exitErr.ExitCode()
		}
		// if exitErr, ok := err.(*exec.ExitError); ok {
		// return exitErr.ExitCode()
		// }
	}

	return 0
}
