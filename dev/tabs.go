package main

// func placeholder() {
// }

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/client"
)

func createStartTab(cli *client.Client) *fyne.Container {
	start_title = canvas.NewText("GoDocker", color.Color(theme.PrimaryColor()))
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

	nrOfServicesEntry := widget.NewEntry()
	// formsContainer := container.NewGridWithColumns(3)
	form := widget.NewForm()
	nameEntry := widget.NewEntry()
	imageOrFileRadio := widget.NewRadioGroup([]string{"Image", "Custom"}, func(s string) {})
	namePathEntry := widget.NewEntry()
	portsCheck := widget.NewCheck("", func(b bool) {})
	hostPortEntry := widget.NewEntry()
	containerPortEntry := widget.NewEntry()

	serviceName := []string{}
	imageOrFile := []string{}
	nameOrPath := []string{}
	bindPorts := []bool{}
	hostPort := []string{}
	containerPort := []string{}

	form.OnSubmit = func() {
		serviceName = append(serviceName, nameEntry.Text)
		imageOrFile = append(imageOrFile, imageOrFileRadio.Selected)
		nameOrPath = append(nameOrPath, namePathEntry.Text)
		bindPorts = append(bindPorts, portsCheck.Checked)
		hostPort = append(hostPort, hostPortEntry.Text)
		containerPort = append(containerPort, containerPortEntry.Text)
	}
	container_compose := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Number of services/ containers:"),
			nrOfServicesEntry,
			widget.NewButton("Create form", func() {
				if nrOfServicesEntry.Text != "" {

					// nrOfServices, err := strconv.Atoi(nrOfServicesEntry.Text)
					// if err != nil {
					// 	log.Fatal(err)
					// }

					// currentNrOfObjects := len(formsContainer.Objects)
					// for i := 1; i <= nrOfServices-currentNrOfObjects; i++ {
					// 	formsContainer.Add(widget.NewForm(
					// 		widget.NewFormItem(
					// 			"Name of the service:",
					// 			widget.NewEntry(),
					// 		),
					// 		widget.NewFormItem(
					// 			"Standard image or custom dockerfile:",
					// 			widget.NewRadioGroup([]string{"Image", "Custom"}, func(s string) {}),
					// 		),
					// 		widget.NewFormItem(
					// 			"Image name/ path to folder:",
					// 			widget.NewEntry(),
					// 		),
					// 		widget.NewFormItem(
					// 			"How many bound ports do you want?",
					// 			widget.NewCheck("", func(b bool) {}),
					// 		),
					// 	))

					// 	formsContainer.Add(widget.NewForm(widget.NewFormItem("Test", widget.NewEntry())))
					// }

					form.Append("Name of the service", nameEntry)
					form.Append("Image or dockerfile", imageOrFileRadio)
					form.Append("Name or path", namePathEntry)
					form.Append("Bind ports?", portsCheck)
					form.Append("Host port:", hostPortEntry)
					form.Append("Container port:", containerPortEntry)

				}
			}),
		),
		form,
		layout.NewSpacer(),
		widget.NewButton("Check", func() {
			fmt.Println(serviceName)
			fmt.Println(imageOrFile)
			fmt.Println(nameOrPath)
			fmt.Println(bindPorts)
			fmt.Println(hostPort)
			fmt.Println(containerPort)
		}),
		layout.NewSpacer(),
		widget.NewButton("Generate", func() {
			fmt.Println(serviceName)
			fmt.Println(imageOrFile)
			fmt.Println(nameOrPath)
			fmt.Println(bindPorts)
			fmt.Println(hostPort)
			fmt.Println(containerPort)
			generateCompose(serviceName, imageOrFile, nameOrPath, bindPorts, hostPort, containerPort)
		}),
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
