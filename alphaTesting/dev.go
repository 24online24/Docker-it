package main

import (
	"context"
	"log"
	"os/exec"
	"strconv"

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

	var showRunningContainersButton *widget.Button = widget.NewButton("Running containers", func() {
		showRunningContainers(w, cli)
	})

	var createComposeFileButton *widget.Button = widget.NewButton("Create Docker Compose file", func() {
		showRunningContainers(w, cli)
	})

	var mainMenu *fyne.Container = container.NewAdaptiveGrid(
		2,
		container.NewVBox(
			container.NewVBox(showRunningContainersButton),
		),
		container.NewVBox(
			container.NewVBox(createComposeFileButton),
		),
	)
	w.SetContent(mainMenu)
	w.ShowAndRun()
}

func showRunningContainers(w fyne.Window, cli *client.Client) {

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	var data [][]string = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES", "TERMINAL"}}
	for _, container := range containers {
		var temp []string = make([]string, 0)
		temp = append(temp, container.ID[:10])
		temp = append(temp, container.Image)
		temp = append(temp, container.Command)
		temp = append(temp, strconv.Itoa(int(container.Created)))
		temp = append(temp, container.Status)

		portString := ""
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
		temp = append(temp, portString)

		temp = append(temp, container.Names...)
		temp = append(temp, "OPEN")
		data = append(data, temp)
	}

	runningContainersTable := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			tab := widget.NewLabel("wide content")
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
			// fmt.Println(data[i.Row][0])
			cmd := exec.Command("cmd", "/c", "start", "cmd", "/c", "docker", "exec", "-ti", data[i.Row][0], "/bin/bash")
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
		runningContainersTable.UnselectAll()
	}

	w.SetContent(runningContainersTable)
}
