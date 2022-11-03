package main

import (
	"context"
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
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

func get_settings() {
	dat, err := os.ReadFile(".settings")
	if err != nil {
		log.Fatal(err)
	}
	str := strings.Split(string(dat), "\n")
	refresh_rate, err = strconv.Atoi(str[0])
	if err != nil {
		log.Fatal(err)
	}
	terminal_setting = str[1]
	docker_path = str[2]
	fmt.Println("Settings have been imported succesfully!")
}

func save_settings() {
	val := fmt.Sprint(refresh_rate) + "\n" + terminal_setting + "\n" + docker_path
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
		time.Sleep(time.Second * time.Duration(refresh_rate))
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

func start_container(status string, data string) {
	var cmd *exec.Cmd
	if strings.Contains(status, "Up") {
		fmt.Println("Open")
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
		fmt.Println("Closed")
		cmd = exec.Command("docker", "start", data)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func showContainers(chContainers chan *widget.Table, cli *client.Client) {
	for {
		var data [][]string = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES", "ACTION"}}

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Limit: 10})
		if err == nil {
			for _, container := range containers {
				var portString string = ""
				for index, port := range container.Ports {
					if len(port.IP) > 0 {
						portString += port.IP + ":"
					}
					if port.PublicPort != 0 {
						portString += strconv.Itoa(int(port.PublicPort)) + "->"
					}
					portString += strconv.Itoa(int(port.PrivatePort)) + "/" + port.Type
					if index < len(container.Ports)-1 {
						portString += ", "
					}
				}

				var actionName string
				if actionName = "Open"; strings.Contains(container.Status, "Up") {
					actionName = "Attach"
				}

				data = append(data, []string{
					container.ID[:10], container.Image, container.Command,
					niceTimeFormat(int(time.Now().Unix() - container.Created)), container.Status,
					portString, container.Names[0], actionName,
				})
			}
		}
		ContainersTable := widget.NewTable(
			func() (int, int) {
				return len(data), len(data[0])
			},
			func() fyne.CanvasObject {
				tab := widget.NewLabel("")
				return tab
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(data[i.Row][i.Col])
			},
		)
		for column := 0; column < len(data[0]); column++ {
			var maxCharLen int = 0
			for row := 0; row < len(data); row++ {
				if len(data[row][column]) > maxCharLen {
					maxCharLen = len(data[row][column])
				}
			}
			ContainersTable.SetColumnWidth(column, 40+float32(maxCharLen)*7)
		}
		ContainersTable.OnSelected = func(i widget.TableCellID) {
			if i.Col == 7 && i.Row > 0 {
				start_container(data[i.Row][4], data[i.Row][0])
			}
			ContainersTable.UnselectAll()
		}
		chContainers <- ContainersTable
		time.Sleep(time.Second * time.Duration(refresh_rate))
	}
}

func showImages(chImages chan *widget.Table, cli *client.Client) {
	for {
		var data [][]string = [][]string{{"REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE"}}

		images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
		if err == nil {
			for _, image := range images {
				data = append(data, []string{
					strings.Split(image.RepoTags[0], ":")[0], strings.Split(image.RepoTags[0], ":")[1],
					strings.Split(image.ID, ":")[1][:12], niceTimeFormat(int(time.Now().Unix() - image.Created)),
					niceSizeFormat(int(image.Size)),
				})
			}
		}
		ImagesTable := widget.NewTable(
			func() (int, int) {
				return len(data), len(data[0])
			},
			func() fyne.CanvasObject {
				tab := widget.NewLabel("")
				return tab
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(data[i.Row][i.Col])
			},
		)
		for column := 0; column < len(data[0]); column++ {
			var maxCharLen int = 0
			for row := 0; row < len(data); row++ {
				if len(data[row][column]) > maxCharLen {
					maxCharLen = len(data[row][column])
				}
			}
			ImagesTable.SetColumnWidth(column, 40+float32(maxCharLen)*7)
		}
		ImagesTable.OnSelected = func(i widget.TableCellID) {
			ImagesTable.UnselectAll()
		}
		time.Sleep(time.Second * time.Duration(refresh_rate))
	}
}

func showVolumes(chVolumes chan *widget.Table, cli *client.Client) {
	for {
		var data [][]string = [][]string{{"DRIVER", "VOLUME NAME"}}

		volumes, err := cli.VolumeList(context.Background(), filters.Args{})
		if err == nil {
			for _, volume := range volumes.Volumes {
				data = append(data, []string{
					volume.Driver, volume.Name,
				})
			}
		}
		VolumesTable := widget.NewTable(
			func() (int, int) {
				return len(data), len(data[0])
			},
			func() fyne.CanvasObject {
				tab := widget.NewLabel("")
				return tab
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(data[i.Row][i.Col])
			},
		)
		for column := 0; column < len(data[0]); column++ {
			var maxCharLen int = 0
			for row := 0; row < len(data); row++ {
				if len(data[row][column]) > maxCharLen {
					maxCharLen = len(data[row][column])
				}
			}
			VolumesTable.SetColumnWidth(column, 40+float32(maxCharLen)*7)
		}
		VolumesTable.OnSelected = func(i widget.TableCellID) {
			VolumesTable.UnselectAll()
		}
		time.Sleep(time.Second * time.Duration(refresh_rate))
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
					refresh_rate, _ = strconv.Atoi(rrate.Text)
					terminal_setting = terminal.Text
					docker_path = docker_e.Text
					save_settings()
				}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
					get_settings()
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
