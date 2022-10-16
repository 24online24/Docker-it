package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")
	myWindow.Resize(fyne.NewSize(720, 480))

	cmd := exec.Command("whoami")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WELCOME %s\n", out)

	content := widget.NewButton("Start Docker Service", func() {

		if runtime.GOOS == "windows" {
			cmd = exec.Command("C:/Program Files/Docker/Docker/Docker Desktop.exe")
		} else if runtime.GOOS == "linux" {
			cmd = exec.Command("systemctl", "start", "docker")
		} else {
			cmd = exec.Command("open", "-a", "Docker")
		}
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started %s", out)
	})

	// check for docker start "ps -C docker -opid=" - ca sa vezi daca e pornit docker sau nu inca, si apoi sa poti sa il opersti(?)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Start", theme.HomeIcon(), content),
		container.NewTabItem("Containers", content),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
