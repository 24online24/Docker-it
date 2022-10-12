package main

import (
	"context"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	a := app.New()
	w := a.NewWindow("GoDocker Containers")
	w.Resize(fyne.NewSize(600, 500))
	var data = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"}}
	for _, container := range containers {
		temp := make([]string, 0)
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
		data = append(data, temp)
		fmt.Println(temp)

	}

	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	w.SetContent(list)
	w.ShowAndRun()
}
