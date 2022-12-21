package main

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Se adaugă la meniul aplicației ferestre care se actualizează constant în funcție de rata de împrospătare a setată.
func updatingTabs(cli *client.Client, tabs *container.AppTabs, function func(ch chan *widget.Table, cli *client.Client), tabName string, index int) {
	tabs.Append(container.NewTabItemWithIcon(tabName, theme.MenuIcon(), widget.NewLabel("")))
	ch := make(chan *widget.Table)
	go function(ch, cli)
	go func() {
		for table := range ch {
			tabs.Items[index].Content = table
		}
	}()
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

func start_container(status string, data string) {
	var cmd *exec.Cmd
	if strings.Contains(status, "Up") {
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
		handleError(err)
	} else {
		cmd = exec.Command("docker", "start", data)
		err := cmd.Run()
		handleError(err)
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
		chImages <- ImagesTable
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
		chVolumes <- VolumesTable
		time.Sleep(time.Second * time.Duration(refresh_rate))
	}
}
