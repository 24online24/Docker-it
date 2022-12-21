package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/docker/docker/client"
)

var env string
var terminal_setting string = ""
var refresh_rate int = 5
var docker_path string = "C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe"
var theme_color string = "dark+"
var start_title *canvas.Text

// Se verifică mediul de rulare: Windows, Linux sau alt sistem nesuportat.
func get_env() {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		env = runtime.GOOS
	} else {
		fmt.Println("Your operating system is not supported by our project. Sorry! D:" + runtime.GOOS)
		os.Exit(0)
	}
}

// Verifică dacă a apărut o eroare. Dacă da, programul se oprește.
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Funcția de unde pornește programul. Aici este creat meniul principal.
func main() {
	get_env()
	get_settings()
	fmt.Println("Running in " + env + " mode...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	handleError(err)

	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("DockerIT")

	tabs := container.NewAppTabs(container.NewTabItemWithIcon("Start", theme.HomeIcon(), createStartTab(cli)))
	updatingTabs(cli, tabs, showContainers, "Containers", 1)
	updatingTabs(cli, tabs, showImages, "Images", 2)
	updatingTabs(cli, tabs, showVolumes, "Volumes", 3)

	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), createComposeTab(cli)))
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), createSettingsTab(cli, a)))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Refresh()
	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
