package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func start_daemon() {
	if env == "windows" {
		cmd := exec.Command(docker_path)
		go cmd.Run()
	} else {
		cmd := exec.Command("systemctl", "start", "docker")
		_ = cmd.Run()
	}
	if !check_daemon() {
		fmt.Println("It hasn't started yet.")
	}
}

func stop_daemon() {
	if check_daemon() {
		if env == "windows" {
			cmd := exec.Command("taskkill", "/im", "Docker Desktop.exe", "/t", "/f")
			cmd.Run()
		} else {
			cmd := exec.Command("systemctl", "stop", "docker*")
			_ = cmd.Run()
		}
		if check_daemon() {
			fmt.Println("It hasn't stopped yet.")
		}
	}
}

func check_daemon() bool {
	cmd2 := exec.Command("docker", "ps")
	out, err := cmd2.CombinedOutput()
	fmt.Println(string(out))
	if strings.Contains(string(out), "error during connect:") ||
		strings.Contains(string(out), "Cannot connect to the Docker daemon") {
		return false
	}
	if err != nil {
		fmt.Println(err)
	}
	return true
}

func isDockerStarted(chDockerStarted chan int) {
	x := 0
	for {
		cmd := exec.Command("docker", "ps")
		out, err := cmd.CombinedOutput()
		if err != nil {
			x = 1
		}
		if strings.Contains(string(out), "error during connect:") ||
			strings.Contains(string(out), "Cannot connect to the Docker daemon") {
			x = 2
		} else {
			x = 3
		}
		time.Sleep(time.Second)
		chDockerStarted <- x
	}
}
