package main

import (
	"context"
	"fmt"
	"image/color"
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
	"github.com/docker/docker/client"
)

func get_env() string {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		return runtime.GOOS
	}
	fmt.Println("Your operating system is not supported by our project. Sorry! D:" + runtime.GOOS)
	os.Exit(0)
	return "nope"
}

func start_daemon(env string) {
	if env == "windows" {
		cmd := exec.Command("powershell", "Start-Process", "'C:\\Program Files\\Docker\\Docker\\resources\\dockerd.exe'", "-WindowStyle", "Hidden")
		_ = cmd.Run()
		if check_daemon(env) == false {
			fmt.Println("BRUH XD")
			os.Exit(0)
		}
	}
}

func check_daemon(env string) bool {
	if env == "windows" {
		cmd2 := exec.Command("powershell", "Get-Process", "dockerd")
		out, err := cmd2.CombinedOutput()
		if strings.Contains(string(out), "Get-Process : Cannot find a process with the name \"dockerd\".") {
			return false
		}
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("dockerd started with success!")
		return true
	}
	return false
}

func showMainMenu(w fyne.Window, cli *client.Client, env string) {
	// Button that starts Docker if it is not started already
	var runDockerButton *widget.Button = widget.NewButton("Start Docker Service", func() {
		if check_daemon(env) == false {
			start_daemon(env)
		} else {
			fmt.Println("daemon already started!")
		}
	})

	// Label showing the current status of Docker
	var dockerd_status *widget.Label = widget.NewLabel("")
	ch := make(chan int)
	go isDockerStarted(ch)
	go func() {
		for running := range ch {
			if running == 3 {
				dockerd_status.SetText("Docker is running! :)")
				fmt.Println("Docker is running! :)")
			} else {
				dockerd_status.SetText("Docker is not running! :(")
				fmt.Println("Docker is not running! :(")
			}
		}
	}()

	// Button that changes the view to that of the running containers
	var showRunningContainersButton *widget.Button = widget.NewButton("Running containers", func() {
		showRunningContainers(w, cli, env)
	})

	var createComposeFileButton *widget.Button = widget.NewButton("Create Docker Compose file", func() {
		showRunningContainers(w, cli, env)
	})

	var mainMenu *fyne.Container = container.NewAdaptiveGrid(
		2,
		container.NewVBox(
			container.NewVBox(runDockerButton, showRunningContainersButton),
		),
		container.NewVBox(
			container.NewVBox(dockerd_status, createComposeFileButton),
		),
	)
	w.SetContent(mainMenu)
}

func showRunningContainers(w fyne.Window, cli *client.Client, env string) {
	go func() {
		var showing bool = true
		for showing {
			containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
			if err != nil {
				panic(err)
			}
			var data [][]string = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES", "TERMINAL"}}
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
				data = append(data, []string{
					container.ID[:10], container.Image, container.Command,
					strconv.Itoa(int(time.Now().Unix() - container.Created)), container.Status,
					portString, container.Names[0], "OPEN",
				})
			}
			runningContainersTable := widget.NewTable(
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
				runningContainersTable.SetColumnWidth(column, 50+float32(maxCharLen)*7)
			}
			runningContainersTable.OnSelected = func(i widget.TableCellID) {
				if i.Col == 7 && i.Row > 0 {
					var cmd *exec.Cmd
					if runtime.GOOS == "windows" {
						cmd = exec.Command("cmd", "/c", "start", "cmd", "/c", "docker", "exec", "-ti", data[i.Row][0], "/bin/bash")
					} else if runtime.GOOS == "linux" {
						testcmd := exec.Command("command", "-v", "gnome-terminal")
						testerr := testcmd.Run()
						if testerr == nil {
							cmd = exec.Command("gnome-terminal", "-e", "docker", "exec", "-ti", data[i.Row][0], "/bin/bash")
						} else {
							testcmd := exec.Command("command", "-v", "konsole")
							testerr := testcmd.Run()
							if testerr == nil {
								cmd = exec.Command("konsole", "-e", "docker", "exec", "-ti", data[i.Row][0], "/bin/bash")
							}
						}
					} else {
						fmt.Println("Bozo")
					}
					err := cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
				} else {
					showing = false
					showMainMenu(w, cli, env)
				}
				runningContainersTable.UnselectAll()
			}
			w.SetContent(runningContainersTable)
			time.Sleep(time.Second)
		}
	}()
}

func isDockerStarted(ch chan int) {
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
		ch <- x
	}
}

func main() {
	env := get_env()
	fmt.Println("Running in " + env + " mode")

	// cli, err := client.NewClientWithOpts(client.FromEnv)
	// if err != nil {
	// 	panic(err)
	// }

	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("GoDocker Containers")

	start_title := canvas.NewText("#RIPBOZO", color.RGBA{255, 0, 0, 3})
	start_title.TextSize = 50

	dockerd_status := widget.NewLabel("")

	start_button := widget.NewButton("Start/Stop", func() {
		if check_daemon(env) == false {
			start_daemon(env)
		} else {
			// TODO add stop function
			fmt.Println("daemon already started!")
		}
	})
	ch := make(chan int)
	go isDockerStarted(ch)
	go func() {
		for running := range ch {
			if running == 3 {
				dockerd_status.SetText("Docker is running! :)")
				fmt.Println("Docker is running! :)")
			} else {
				dockerd_status.SetText("Docker is not running! :(")
				fmt.Println("Docker is not running! :(")
			}
		}
	}()

	container_start := container.NewVBox(
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
		container.New(layout.NewCenterLayout(), widget.NewLabel("Guide here!\nAsta se mai centreaza dupa annyway!")),
		layout.NewSpacer(),
	)

	tabs := container.NewAppTabs(container.NewTabItemWithIcon("Start", theme.HomeIcon(), container_start))
	tabs.Append(container.NewTabItemWithIcon("Containers", theme.MenuIcon(), container.NewVBox(dockerd_status)))
	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), widget.NewLabel("Compose Docker files goes here!")))
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), widget.NewLabel("Settings goes here!")))
	tabs.SetTabLocation(container.TabLocationTop)

	// showMainMenu(w, cli, env)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(720, 480))
	w.ShowAndRun()
}
