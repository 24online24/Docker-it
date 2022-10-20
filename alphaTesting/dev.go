package main

import (
	"context"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("GoDocker Containers")
	w.Resize(fyne.NewSize(1600, 900))

	showMainMenu(w, cli)
	w.ShowAndRun()
}

func showMainMenu(w fyne.Window, cli *client.Client) {
	// Button that starts Docker if it is not started already
	var runDockerButton *widget.Button = widget.NewButton("Start Docker Service", func() {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("C:/Program Files/Docker/Docker/Docker Desktop.exe")
		} else if runtime.GOOS == "linux" {
			cmd = exec.Command("systemctl", "start", "docker")
		} else {
			cmd = exec.Command("open", "-a", "Docker")
		}
		err := cmd.Run()
		// fmt.Printf("%s\n", out)
		if err != nil {
			log.Fatal(err)
			// fmt.Printf("Alrdy started! %s\n", err)
		}
	})

	// Label showing the current status of Docker
	var isDockerStartedLabel *widget.Label = widget.NewLabel("")
	ch := make(chan int)
	go isDockerStarted(ch)
	go func() {
		for running := range ch {
			if running == 3 {
				isDockerStartedLabel.SetText("Docker is running! :)")
				// fmt.Println("Docker is running! :)")
			} else {
				isDockerStartedLabel.SetText("Docker is not running! :(")
				// fmt.Println("Docker is not running! :(")
			}
		}
	}()

	// Changes the view to that of the running containers
	var showRunningContainersButton *widget.Button = widget.NewButton("Running containers", func() {
		showRunningContainers(w, cli)
	})

	// Changes the view to the docker-compose.yml creation windows
	var createComposeFileButton *widget.Button = widget.NewButton("Create Docker Compose file", func() {
		showRunningContainers(w, cli)
	})

	// Menu organized in 2 columns
	var mainMenu *fyne.Container = container.NewAdaptiveGrid(
		2,
		container.NewVBox(
			container.NewVBox(runDockerButton, showRunningContainersButton),
		),
		container.NewVBox(
			container.NewVBox(isDockerStartedLabel, createComposeFileButton),
		),
	)
	w.SetContent(mainMenu)
}

// View that shows the currently running Docker containers. It is refreshed once
// every second and allows the user to open a terminal directly connected
// to a container.
func showRunningContainers(w fyne.Window, cli *client.Client) {
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
					}
					// else {
					// 	fmt.Println("Bozo")
					// }
					err := cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
				} else {
					showing = false
					showMainMenu(w, cli)
				}
				runningContainersTable.UnselectAll()
			}
			w.SetContent(runningContainersTable)
			time.Sleep(time.Second)
		}
	}()
}

// Background check to see if Docker is started or not
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
		ch <- x
		time.Sleep(time.Second)
	}
}
