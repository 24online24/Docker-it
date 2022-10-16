package main

import (
	"log"
	"os/exec"
)

func main() {
	// cmd := exec.Command("cmd", "/c", "start", "cmd")
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/c", "docker", "exec", "-ti", "2f2e499d758d", "/bin/bash")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
