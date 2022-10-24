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
		// cmd := exec.Command("powershell", "Start-Process", "'C:\\Program Files\\Docker\\Docker\\resources\\dockerd.exe'", "-WindowStyle", "Hidden")
		cmd := exec.Command("C:/Program Files/Docker/Docker/Docker Desktop.exe")
		go cmd.Run()
	} else {
		cmd := exec.Command("systemctl", "start", "docker")
		_ = cmd.Run()
	}
	if !check_daemon() {
		fmt.Println("BRUH it didnt start XD")
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

func showRunningContainers(chRunningContainers chan *widget.Table, cli *client.Client) {
	go func() {
		for {
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
					}
					// else {
					// 	fmt.Println("Bozo")
					// }
					err := cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
				}
				runningContainersTable.UnselectAll()
			}
			chRunningContainers <- runningContainersTable
			time.Sleep(time.Second)
		}
	}()
}

func main() {
	env := get_env()
	fmt.Println("Running in " + env + " mode...")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err)
	}

	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("GoDocker Containers")

	start_title := canvas.NewText("#RIPBOZO", color.RGBA{255, 0, 0, 3})
	start_title.TextSize = 50

	dockerd_status := widget.NewLabel("")

	start_button := widget.NewButton("Start/Stop", func() {
		if !check_daemon() {
			start_daemon(env)
		} else {
			// TODO add stop function
			fmt.Println("daemon already started!")
		}
	})

	chDockerStarted := make(chan int)
	go isDockerStarted(chDockerStarted)
	go func() {
		for running := range chDockerStarted {
			if running == 3 {
				dockerd_status.SetText("Docker is running! :)")
				// fmt.Println("Docker is running! :)")
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

	if check_daemon() {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			fmt.Println(err.Error())
		}

		var data [][]string = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES", "TERMINAL"}}
		for _, container := range containers {
			var portString string = ""
			fmt.Printf("THIS IS CONTAINER" + container.Names[0])
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

		container_4containers := widget.NewTable(
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
			container_4containers.SetColumnWidth(column, 50+float32(maxCharLen)*7)
		}

		container_4containers.OnSelected = func(i widget.TableCellID) {
			if i.Col == 7 && i.Row > 0 {
				// TODO could be integrated in a function
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
				}
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
			container_4containers.UnselectAll()
			time.Sleep(time.Second)
		}
		tabs.Append(container.NewTabItemWithIcon("Containers", theme.MenuIcon(), container_4containers))
	} else {
		tabs.Append(container.NewTabItemWithIcon("Containers", theme.MenuIcon(), widget.NewLabel("XDD")))
	}

	chRunningContainers := make(chan *widget.Table)
	go showRunningContainers(chRunningContainers, cli)
	go func() {
		for table := range chRunningContainers {
			tabs.Items[1].Content = table
		}
	}()

	// container_compose := container.NewWithoutLayout()
	// container_settings := container.NewWithoutLayout()

	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), widget.NewLabel("Compose Docker files goes here!")))
	tabs.Append(container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), widget.NewLabel("Settings goes here!")))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Refresh()
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
