package main

import (
	"context"
	"fmt"
	"strconv"

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

	// a := app.New()
	// w := a.NewWindow("GoDocker Containers")
	// w.Resize(fyne.NewSize(600, 500))
	// var data = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"}}
	for _, container := range containers {
		// temp := make([]string, 0)
		// temp = append(temp, container.ID[:10])
		// temp = append(temp, container.Image)
		// temp = append(temp, container.Command)
		// temp = append(temp, strconv.Itoa(int(container.Created)))
		// temp = append(temp, container.Status)
		// temp = append(temp, container.Ports)

		// temp = append(temp, container.Names...)
		// data = append(data, temp)
		// fmt.Println(container.ID[:10] + " " + container.Image)
		for _, port := range container.Ports {
			if len(port.IP) > 0 {
				fmt.Print(port.IP + ":")
			}
			if port.PublicPort != 0 {
				fmt.Print(strconv.Itoa(int(port.PublicPort)) + "->")
			}
			fmt.Print(strconv.Itoa(int(port.PrivatePort)))
			fmt.Print("/" + port.Type)
			fmt.Print("\t\t")
		}
		// fmt.Print(container.Ports)
		fmt.Println()
	}

	// list := widget.NewTable(
	// 	func() (int, int) {
	// 		return len(data), len(data[0])
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return widget.NewLabel("wide content")
	// 	},
	// 	func(i widget.TableCellID, o fyne.CanvasObject) {
	// 		o.(*widget.Label).SetText(data[i.Row][i.Col])
	// 	})

	// w.SetContent(list)
	// w.ShowAndRun()
}
