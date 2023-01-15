package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// returns a table with all the containers and their information
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

// returns a table with all the images and their information
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
		chImages <- ImagesTable
		time.Sleep(time.Second * time.Duration(refresh_rate))
	}
}

// returns a table with all the volumes and their information
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
		chVolumes <- VolumesTable
		time.Sleep(time.Second * time.Duration(refresh_rate))
	}
}
