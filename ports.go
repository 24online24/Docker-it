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

	var a fyne.App = app.New()
	var w fyne.Window = a.NewWindow("GoDocker Containers")
	w.Resize(fyne.NewSize(1600, 900))
	// w.SetFullScreen(true)
	// var data = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND"}}
	var data [][]string = [][]string{{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"}}
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
				// if index%2 == 1 {
				// 	portString += "\n"
				// }
			}
		}
		temp = append(temp, portString)

		temp = append(temp, container.Names...)
		data = append(data, temp)
		// fmt.Println(temp)

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
		})

	/*
		For every column we check every row for the longest entry (string).
		If there are multiple lines in an entry, the length of the first line is taken into account.
		We then format the column according to the longest entry.
	*/
	for column := 0; column < len(data[0]); column++ {
		// var lineLenMax int = 0
		var maxCharLen int = 0
		// var lines int = 0
		for row := 0; row < len(data); row++ {
			if len(data[row][column]) > maxCharLen {
				maxCharLen = len(data[row][column])
				// var lineLen int = 0
				// var newLineIndex = strings.Index(data[row][column], "\n")
				// if newLineIndex > 0 {
				// 	lineLen = newLineIndex
				// } else {
				// 	lineLen = len(data[row][column])
				// }
				// if lineLen > lineLenMax {
				// 	lineLenMax = lineLen
				// }
			}
		}
		runningContainersTable.SetColumnWidth(column, 50+float32(maxCharLen)*7)
	}
	w.SetContent(runningContainersTable)

	button := widget.NewButton("Running containers", func() {
		fmt.Println("tapped")
		w.SetContent(runningContainersTable)
	})
	w.SetContent(button)
	w.ShowAndRun()
}
