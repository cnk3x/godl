package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func runShell(script string) error {
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	log.Println(cmd.String())
	return cmd.Run()
}

func linkDo(source, target string) error {
	return runShell(fmt.Sprintf("ln -f -s %s %s", source, target))
}
