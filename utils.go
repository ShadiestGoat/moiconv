package main

import (
	"os"
	"os/exec"
)

func requireBin(name string) {
	_, err := exec.LookPath("exiftool")
	if err != nil {
		panic("Couldn't find " + name + "; is it installed?")
	}
}

func runCmd(args ...string) error {
	cmd := exec.Command(
		args[0],
		args[1:]...
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = nil

	return cmd.Run()
}
