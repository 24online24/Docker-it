package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/client"
)

var env string                                                                   // operating system
var terminal_setting string = ""                                                 // for custom terminal opening
var refresh_rate int = 5                                                         // refresh rate in seconds
var docker_path string = "C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe" // path to docker desktop (windows only)
var theme_color string = "dark+"                                                 // theme color
var start_title *canvas.Text                                                     // title of start tab

// get operating system
func get_env() {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		env = runtime.GOOS
	} else {
		fmt.Println("Your operating system is not supported by our project. Sorry! D:" + runtime.GOOS)
		os.Exit(0)
	}
}

// function to open container terminal
func start_container(status string, data string) {
	var cmd *exec.Cmd
	if strings.Contains(status, "Up") {
		if runtime.GOOS == "windows" {
			if terminal_setting == "" {
				cmd = exec.Command("cmd", "/c", "start", "cmd", "/c", "docker", "exec", "-ti", data, "/bin/bash")
			} else {
				rest := [5]string{"docker", "exec", "-ti", data, "/bin/bash"}
				cmd_line := strings.Split(terminal_setting, " ")
				for i := 0; i < 5; i++ {
					cmd_line = append(cmd_line, rest[i])
				}
				cmd = exec.Command("powershell.exe", cmd_line...)

			}
		} else if runtime.GOOS == "linux" {
			testcmd := exec.Command("command", "-v", "gnome-terminal")
			testerr := testcmd.Run()
			if testerr == nil {
				cmd = exec.Command("gnome-terminal", "-e", "docker", "exec", "-ti", data, "/bin/bash")
			} else {
				testcmd := exec.Command("command", "-v", "konsole")
				testerr := testcmd.Run()
				if testerr == nil {
					cmd = exec.Command("konsole", "-e", "docker", "exec", "-ti", data, "/bin/bash")
				}
			}
		}
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cmd = exec.Command("docker", "start", data)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	get_env()
	get_settings() // get settings from .settings file
	fmt.Println("Running in " + env + " mode...")
	cli, err := client.NewClientWithOpts(client.FromEnv) // docker sdk client
	if err != nil {
		log.Fatal(err)
	}

	// initialize fyne app and window
	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("DockerIT")

	// create start tab
	tabs := container.NewAppTabs(container.NewTabItemWithIcon("Start", theme.HomeIcon(), createStartTab(cli)))

	// create containers tab
	tabs.Append(container.NewTabItemWithIcon("Containers", theme.MenuIcon(), widget.NewLabel("")))
	chContainers := make(chan *widget.Table)
	go showContainers(chContainers, cli) // goroutine to show containers
	go func() {                          // goroutine to update containers
		for table := range chContainers {
			tabs.Items[1].Content = table
		}
	}()

	// create images tab
	tabs.Append(container.NewTabItemWithIcon("Images", theme.MenuIcon(), widget.NewLabel("")))
	chImages := make(chan *widget.Table)
	go showImages(chImages, cli) // goroutine to show images
	go func() {                  // goroutine to update images
		for table := range chImages {
			tabs.Items[2].Content = table
		}
	}()

	// create volumes tab
	tabs.Append(container.NewTabItemWithIcon("Volumes", theme.MenuIcon(), widget.NewLabel("")))
	chVolumes := make(chan *widget.Table)
	go showVolumes(chVolumes, cli) // goroutine to show volumes
	go func() {                    // goroutine to update volumes
		for table := range chVolumes {
			tabs.Items[3].Content = table
		}
	}()

	// theme select widget
	theme_options := []string{"dark", "dark+", "light"}
	theme_select := widget.NewSelect(theme_options, func(s string) {
		switch s {
		case "dark":
			a.Settings().SetTheme(theme.DefaultTheme())
			theme_color = s
		case "dark+":
			a.Settings().SetTheme(&dark_plus{})
			theme_color = s
		case "light":
			a.Settings().SetTheme(&white_theme{})
			theme_color = s
		}
	})

	// selected theme from .settings file or default
	theme_select.SetSelected(theme_color)

	// termianl input widget
	terminal := widget.NewEntry()
	if terminal_setting != "" {
		terminal.Text = terminal_setting
	} else {
		terminal.SetPlaceHolder("Terminal path/executable goes here...")
	}

	// docker path input widget (windows only)
	docker_e := widget.NewEntry()
	if env == "linux" {
		docker_e.SetPlaceHolder("This setting is for windows only :D")
		docker_e.Disable()
	} else {
		if docker_path != "" {
			docker_e.Text = docker_path
		} else {
			docker_e.Text = "C:/Program Files/Docker/Docker/Docker Desktop.exe"
		}
	}

	// refresh rate input widget
	rrate := widget.NewEntry()
	rrate.SetPlaceHolder("Number in seconds")
	if refresh_rate == 0 {
		rrate.Text = "1"
	} else {
		rrate.Text = fmt.Sprint(refresh_rate)
	}
	rrate.Validate()

	// compose tab
	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), createComposeTab(cli)))
	// settings tab
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), createSettingsTab(cli, theme_select, terminal, docker_e, rrate)))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Refresh()
	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
