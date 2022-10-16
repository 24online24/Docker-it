package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func main() {
	cmd := exec.Command("whoami")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WELCOME %s\n", out)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("C:/Program Files/Docker/Docker/Docker Desktop.exe")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("systemctl", "start", "docker")
	} else {
		cmd = exec.Command("open", "-a", "Docker")
	}
	out2, err2 := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err2)
	}
	fmt.Printf("%s\n", out2)
}
