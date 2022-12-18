package main

func placeholder() {
}

// import (
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/layout"
// 	"fyne.io/fyne/v2/theme"
// 	"fyne.io/fyne/v2/widget"
// 	"github.com/docker/docker/client"
// )

// func createTabs(cli *client.Client) {
// 	container_start := container.NewVBox(
// 		layout.NewSpacer(),
// 		container.New(layout.NewCenterLayout(), start_title),
// 		layout.NewSpacer(),
// 		container.New(layout.NewGridLayoutWithColumns(4),
// 			layout.NewSpacer(),
// 			widget.NewLabel("Docker daemon status:"),
// 			dockerd_status,
// 			layout.NewSpacer(),
// 		),
// 		container.New(layout.NewGridLayoutWithColumns(4),
// 			layout.NewSpacer(),
// 			widget.NewLabel("Start/Stop Daemon:"),
// 			start_button,
// 			layout.NewSpacer(),
// 		),
// 		layout.NewSpacer(),
// 		// container.New(layout.NewCenterLayout(), widget.NewLabel("Help")),
// 		layout.NewSpacer(),
// 	)

// 	tabs := container.NewAppTabs(container.NewTabItemWithIcon("Start", theme.HomeIcon(), container_start))

// 	tabs.Append(container.NewTabItemWithIcon("Containers", theme.MenuIcon(), widget.NewLabel("")))
// 	chContainers := make(chan *widget.Table)
// 	go showContainers(chContainers, cli)
// 	go func() {
// 		for table := range chContainers {
// 			tabs.Items[1].Content = table
// 		}
// 	}()

// 	tabs.Append(container.NewTabItemWithIcon("Images", theme.MenuIcon(), widget.NewLabel("")))
// 	chImages := make(chan *widget.Table)
// 	go showImages(chImages, cli)
// 	go func() {
// 		for table := range chImages {
// 			tabs.Items[2].Content = table
// 		}
// 	}()

// 	tabs.Append(container.NewTabItemWithIcon("Volumes", theme.MenuIcon(), widget.NewLabel("")))
// 	chVolumes := make(chan *widget.Table)
// 	go showVolumes(chVolumes, cli)
// 	go func() {
// 		for table := range chVolumes {
// 			tabs.Items[3].Content = table
// 		}
// 	}()

// 	// container_compose := container.NewWithoutLayout()
// 	tabs.Append(container.NewTabItemWithIcon("Compose", theme.FileIcon(), widget.NewLabel("Compose Docker files goes here!")))

// }
