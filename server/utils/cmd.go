package utils

import (
	"bytes"
	"errors"
	"os/exec"
)

func RunCmd(cmd *exec.Cmd) (bytes.Buffer, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out, errors.New(stderr.String())
	}

	return out, nil
}
