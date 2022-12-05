package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

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

var cli *client.Client

func get_env() {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		env = runtime.GOOS
	} else {
		fmt.Println("Your operating system is not supported by our project. Sorry! D:" + runtime.GOOS)
		os.Exit(0)
	}
}

func get_settings() {
	dat, err := os.ReadFile(".settings")
	if err != nil {
		log.Fatal(err)
	}
	str := strings.Split(string(dat), "\n")
	refresh_rate, err = strconv.Atoi(strings.Trim(str[0], "\r"))
	if err != nil {
		log.Fatal(err)
	}
	terminal_setting = strings.Trim(str[1], "\r")
	if env == "windows" {
		docker_path = strings.Trim(str[2], "\r")
	}
	fmt.Println("Settings have been imported succesfully!")
}

func save_settings() {
	val := fmt.Sprint(refresh_rate) + "\n" + terminal_setting + "\n"
	if env == "windows" {
		val += docker_path
	}
	data := []byte(val)

	err := ioutil.WriteFile(".settings", data, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Settings have been saved succesfully!")
}

func start_daemon() {
	if env == "windows" {
		cmd := exec.Command(docker_path)
		go cmd.Run()
	} else {
		cmd := exec.Command("systemctl", "start", "docker")
		_ = cmd.Run()
	}
	if !check_daemon() {
		fmt.Println("BRUH it didnt start XD yet...")
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
			fmt.Println("BRUH it didnt stop XD")
		}
	}
}

func check_daemon() bool {
	cmd2 := exec.Command("docker", "ps")
	out, err := cmd2.CombinedOutput()
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

func niceTimeFormat(seconds int) string {
	var time string = ""
	if seconds > 60 {
		var minutes int = seconds / 60
		seconds %= 60
		if minutes > 60 {
			var hours int = minutes / 60
			minutes = minutes % 60
			if hours > 24 {
				var days int = hours / 24
				hours = hours % 24
				time = strconv.Itoa(days) + "d "
			}
			time = time + strconv.Itoa(hours) + "h "
		}
		time = time + strconv.Itoa(minutes) + "m "
	}
	time = time + strconv.Itoa(seconds) + "s ago"
	return time
}

func niceSizeFormat(bytes int) string {
	var size string = ""
	switch {
	case bytes > 1024*1024*1024:
		size = strconv.Itoa(bytes/1024/1024/1024) + " GiB"
	case bytes > 1024*1024:
		size = strconv.Itoa(bytes/1024/1024) + " MiB"
	case bytes > 1024:
		size = strconv.Itoa(bytes/1024) + " KiB"
	default:
		size = strconv.Itoa(bytes) + " Bytes"
	}
	return size
}

func main() {
	get_env()
	get_settings()
	fmt.Println("Running in " + env + " mode...")
	var errcli error
	cli, errcli = client.NewClientWithOpts(client.FromEnv)
	if errcli != nil {
		log.Fatal(errcli)
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

	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), widget.NewButton("Compose!", clicompose)))
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container_settings))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Refresh()
	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
