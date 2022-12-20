package main

import (
	"fmt"
	"image/color"
	"reflect"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/client"
)

type containerInfo struct {
	serviceName   string
	imageOrFile   string
	nameOrPath    string
	bindPorts     bool
	hostPort      string
	containerPort string
}

func createStartTab(cli *client.Client) *fyne.Container {
	start_title = canvas.NewText("DockerIT", color.Color(theme.PrimaryColor()))
	start_title.TextSize = 50

	dockerd_status := widget.NewLabel("")

	start_button := widget.NewButton("Start/Stop", func() {
		if !check_daemon() {
			start_daemon()
		} else {
			stop_daemon()
		}
	})

	start_button.Importance = widget.HighImportance
	chDockerStarted := make(chan int)
	go isDockerStarted(chDockerStarted)
	go func() {
		for running := range chDockerStarted {
			if running == 3 {
				dockerd_status.SetText("Docker is running! :)")
			} else {
				dockerd_status.SetText("Docker is not running! :(")
			}
		}
	}()

	container_start := container.NewVBox(
		layout.NewSpacer(),
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
		// container.New(layout.NewCenterLayout(), widget.NewLabel("Help")),
		layout.NewSpacer(),
	)

	return container_start
}

func createComposeTab(cli *client.Client) *fyne.Container {

	nameEntry := widget.NewEntry()
	imageOrFileRadio := widget.NewRadioGroup([]string{"Image", "Custom"}, func(s string) {})
	namePathEntry := widget.NewEntry()
	hostPortEntry := widget.NewEntry()
	containerPortEntry := widget.NewEntry()
	portsCheck := widget.NewCheck("", func(b bool) {
		if b {
			hostPortEntry.Enable()
			containerPortEntry.Enable()
		} else {
			hostPortEntry.Disable()
			containerPortEntry.Disable()
		}
	})
	portsCheck.SetChecked(true)
	form := widget.NewForm(
		widget.NewFormItem("Name of the service", nameEntry),
		widget.NewFormItem("Image or dockerfile", imageOrFileRadio),
		widget.NewFormItem("Name or path", namePathEntry),
		widget.NewFormItem("Bind ports?", portsCheck),
		widget.NewFormItem("Host port:", hostPortEntry),
		widget.NewFormItem("Container port:", containerPortEntry),
	)

	containerList := []containerInfo{}
	form.OnSubmit = func() {

		currentContainer := containerInfo{
			serviceName:   nameEntry.Text,
			imageOrFile:   imageOrFileRadio.Selected,
			nameOrPath:    namePathEntry.Text,
			bindPorts:     portsCheck.Checked,
			hostPort:      hostPortEntry.Text,
			containerPort: containerPortEntry.Text,
		}

		if !currentContainer.bindPorts {
			currentContainer.hostPort = "-"
			currentContainer.containerPort = "-"
		}
		containerList = append(containerList, currentContainer)
		nameEntry.SetText("")
		imageOrFileRadio.SetSelected("")
		namePathEntry.SetText("")
		portsCheck.SetChecked(true)
		hostPortEntry.SetText("")
		containerPortEntry.SetText("")
	}

	containerForChecking := container.NewVBox()
	verticalScrollingContainer := container.NewVScroll(containerForChecking)

	checkCompose := func() {
		for {
			nrShown := len(containerForChecking.Objects) / reflect.TypeOf(containerInfo{}).NumField()
			nrContainers := len(containerList)
			for index := nrShown; index < nrContainers; index++ {
				containerForChecking.Add(widget.NewLabel("Service name: " + containerList[index].serviceName))
				containerForChecking.Add(widget.NewLabel("Image/ custom file: " + containerList[index].imageOrFile))
				containerForChecking.Add(widget.NewLabel("Image name/ file path: " + containerList[index].nameOrPath))

				if containerList[index].bindPorts {
					containerForChecking.Add(widget.NewLabel("Bind ports: true"))
					containerForChecking.Add(widget.NewLabel("Host port: " + containerList[index].hostPort))
					containerForChecking.Add(widget.NewLabel("Container port: " + containerList[index].containerPort))
				} else {
					containerForChecking.Add(widget.NewLabel("Bind ports: false"))
				}
				containerForChecking.Add(canvas.NewLine(color.RGBA{128, 128, 128, 255}))
			}
		}
	}
	go checkCompose()

	container_compose := container.NewGridWithColumns(
		2,
		container.NewVBox(
			form,
			layout.NewSpacer(),
			widget.NewButton("Generate", func() {
				generateCompose(containerList)
			}),
		),
		verticalScrollingContainer,
	)

	return container_compose
}

func createSettingsTab(cli *client.Client, theme_select *widget.Select, terminal *widget.Entry, docker_e *widget.Entry, rrate *widget.Entry) *fyne.Container {
	container_settings := container.NewHBox(
		container.NewVBox(
			layout.NewSpacer(),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Select your favourite theme:"),
				theme_select,
				layout.NewSpacer(),
			),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Terminal path to executable:"),
				terminal,
				layout.NewSpacer(),
			),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Docker Desktop.exe path (Windows only):"),
				docker_e,
				layout.NewSpacer(),
			),
			container.New(layout.NewGridLayoutWithColumns(4),
				layout.NewSpacer(),
				widget.NewLabel("Refresh rate for containers in seconds:"),
				rrate,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
			container.NewHBox(
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
					refresh_rate, _ = strconv.Atoi(rrate.Text)
					terminal_setting = terminal.Text
					if env == "windows" {
						docker_path = docker_e.Text
					}
					save_settings()
				}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
					get_settings()
					rrate.Text = fmt.Sprint(refresh_rate)
					terminal.Text = terminal_setting
					if env == "windows" {
						docker_e.Text = docker_path
						docker_e.Refresh()
					}
					rrate.Refresh()
					terminal.Refresh()

				}),
				layout.NewSpacer(),
				layout.NewSpacer(),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		))
	return container_settings
}
