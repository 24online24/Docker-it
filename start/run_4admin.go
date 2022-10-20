package main

import (
	"fmt"
	"os"
	"os/exec"
)

func this() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("Not started in Admin mode. Requesting permissions...")
		return false
	}
	fmt.Println("Admin permissions available. Continuing...")
	return true
}

func main() {
	if this() == false {
		_ = exec.Command("powershell", "Start-Process", "powershell.exe", "-ArgumentList", "\"-noexit cd C:\\Users\\Catalin\\go\\src\\start\\; go run dev.go\"", "-Verb", "RunAs").Run()
	}
}
