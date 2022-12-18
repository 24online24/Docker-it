package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/client"
)

var env string
var terminal_setting string = ""
var refresh_rate int = 1
var docker_path string = ""

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
	var w fyne.Window = a.NewWindow("GoDocker")

	start_title := canvas.NewText("GoDocker", color.RGBA{0, 183, 237, 3})
	start_title.TextSize = 50

	dockerd_status := widget.NewLabel("")

	start_button := widget.NewButton("Start/Stop", func() {
		if !check_daemon() {
			start_daemon()
		} else {
			stop_daemon()
		}
	})

	chDockerStarted := make(chan int)
	go isDockerStarted(chDockerStarted)
	go func() {
		for running := range chDockerStarted {
			if running == 3 {
				dockerd_status.SetText("Docker is running! :)")
			} else {
				dockerd_status.SetText("Docker is not running! :(")
			}
		}
	}()

	// add filter menus, stop container option
	container_start := container.NewVBox(
		layout.NewSpacer(),
		container.New(layout.NewCenterLayout(), start_title),
		layout.NewSpacer(),
		container.New(layout.NewGridLayoutWithColumns(4),
			layout.NewSpacer(),
			widget.NewLabel("Docker daemon status:"),
			dockerd_status,
			layout.NewSpacer(),
		),
		container.New(layout.NewGridLayoutWithColumns(4),
			layout.NewSpacer(),
			widget.NewLabel("Start/Stop Daemon:"),
			start_button,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		// container.New(layout.NewCenterLayout(), widget.NewLabel("Help")),
		layout.NewSpacer(),
	)

	tabs := container.NewAppTabs(container.NewTabItemWithIcon("Start", theme.HomeIcon(), container_start))

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
		case "dark+":
			a.Settings().SetTheme(&darker_than_my_soul{})
		case "light":
			a.Settings().SetTheme(&white_theme{})
		}
	})

	// TODO set to prev selected theme
	theme_select.SetSelected("dark")

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
	rrate.SetPlaceHolder("Number between 1s to 5m...")
	if refresh_rate == 0 {
		rrate.Text = "1"
	} else {
		rrate.Text = fmt.Sprint(refresh_rate)
	}
	rrate.Validate()

	container_settings := container.NewHBox(
		container.NewVBox(
			layout.NewSpacer(),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Select your favourite theme:"),
				theme_select,
				layout.NewSpacer(),
			),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Terminal path to executable:"),
				terminal,
				layout.NewSpacer(),
			),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Docker Desktop.exe path (Windows only):"),
				docker_e,
				layout.NewSpacer(),
			),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Refresh rate for containers in seconds:"),
				rrate,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
			container.NewHBox(
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
					// TODO add theme here aswell
					refresh_rate, _ = strconv.Atoi(rrate.Text)
					terminal_setting = terminal.Text
					if env == "windows" {
						docker_path = docker_e.Text
					}
					save_settings()
				}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
					get_settings()
					rrate.Text = fmt.Sprint(refresh_rate)
					terminal.Text = terminal_setting
					if env == "windows" {
						docker_e.Text = docker_path
						docker_e.Refresh()
					}
					rrate.Refresh()
					terminal.Refresh()

				}),
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		))

	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), widget.NewLabel("Compose Docker files goes here!")))
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container_settings))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Refresh()
	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
