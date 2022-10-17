package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func is_docker_started() {
	for {
		cmd := exec.Command("docker", "stats")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		if strings.Contains(string(out), "error during connect:") {
			fmt.Println("DOCKER NOT STARTED (yet)")
		} else {
			fmt.Println("DOCKER STARTED")
			break
		}

	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")

	go is_docker_started()

	cmd := exec.Command("whoami")
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Printf("Nu stiu cum te cheama %s!", err)
	}

	content := widget.NewButton("Start Docker Service", func() {
		go is_docker_started()
		if runtime.GOOS == "windows" {
			cmd = exec.Command("C:/Program Files/Docker/Docker/Docker Desktop.exe")
		} else if runtime.GOOS == "linux" {
			cmd = exec.Command("systemctl", "start", "docker")
		} else {
			cmd = exec.Command("open", "-a", "Docker")
		}
		out, err = cmd.CombinedOutput()
		fmt.Printf("%s\n", out)
		if err != nil {
			fmt.Printf("Alrdy started! %s\n", err)
		}
	})

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Start", theme.HomeIcon(), content),
		container.NewTabItem("Containers", widget.NewLabel("XDD")),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(720, 480))
	myWindow.ShowAndRun()
}
