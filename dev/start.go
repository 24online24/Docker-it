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

var env string
var terminal_setting string = ""
var refresh_rate int = 5
var docker_path string = "C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe"
var theme_color string = "dark+"
var start_title *canvas.Text

func get_env() {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		env = runtime.GOOS
	} else {
		fmt.Println("Your operating system is not supported by our project. Sorry! D:" + runtime.GOOS)
		os.Exit(0)
	}
}

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
	get_settings()
	fmt.Println("Running in " + env + " mode...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("DockerIT")

	tabs := container.NewAppTabs(container.NewTabItemWithIcon("Start", theme.HomeIcon(), createStartTab(cli)))

	tabs.Append(container.NewTabItemWithIcon("Containers", theme.MenuIcon(), widget.NewLabel("")))
	chContainers := make(chan *widget.Table)
	go showContainers(chContainers, cli)
	go func() {
		for table := range chContainers {
			tabs.Items[1].Content = table
		}
	}()

	tabs.Append(container.NewTabItemWithIcon("Images", theme.MenuIcon(), widget.NewLabel("")))
	chImages := make(chan *widget.Table)
	go showImages(chImages, cli)
	go func() {
		for table := range chImages {
			tabs.Items[2].Content = table
		}
	}()

	tabs.Append(container.NewTabItemWithIcon("Volumes", theme.MenuIcon(), widget.NewLabel("")))
	chVolumes := make(chan *widget.Table)
	go showVolumes(chVolumes, cli)
	go func() {
		for table := range chVolumes {
			tabs.Items[3].Content = table
		}
	}()

	// container_compose := container.NewWithoutLayout()

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
			a.Settings().SetTheme(&myTheme{})
			theme_color = s
		}
		w.Content().Refresh()
	})

	theme_select.SetSelected(theme_color)

	terminal := widget.NewEntry()
	if terminal_setting != "" {
		terminal.Text = terminal_setting
	} else {
		terminal.SetPlaceHolder("Terminal path/executable goes here...")
	}
	docker_e := widget.NewEntry()
	if env == "linux" {
		docker_e.SetPlaceHolder("This setting is for windows only :D")
	} else {
		if docker_path != "" {
			docker_e.Text = docker_path
		} else {
			docker_e.Text = "C:/Program Files/Docker/Docker/Docker Desktop.exe"
		}
	}
	rrate := widget.NewEntry()
	rrate.SetPlaceHolder("Number in seconds")
	if refresh_rate == 0 {
		rrate.Text = "1"
	} else {
		rrate.Text = fmt.Sprint(refresh_rate)
	}
	rrate.Validate()

	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), createComposeTab(cli)))
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), createSettingsTab(cli, theme_select, terminal, docker_e, rrate)))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Refresh()
	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	go tabs.Refresh()
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
