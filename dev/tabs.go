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
	start_title := canvas.NewText("GoDocker", color.RGBA{0, 183, 237, 3})
	start_title.TextSize = 50

	dockerd_status := widget.NewLabel("")

	start_button := widget.NewButton("Start/Stop", func() {
		if !check_daemon() {
			start_daemon()
		} else {
			stop_daemon()
		}
	})

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

	// nrOfServices := 0
	// nrOfServicesEntry := widget.NewEntry()
	// nrOfServicesEntry.SetPlaceHolder("Number between 1s to 5m...")
	// if nrOfServices == 0 {
	// 	nrOfServicesEntry.Text = "1"
	// } else {
	// 	nrOfServicesEntry.Text = fmt.Sprint(nrOfServices)
	// }
	// nrOfServicesEntry.Validate()

	// formsContainer := container.NewGridWithColumns(3)

	type form struct {
		content widget.Form
	}
	forms := []form{}

	forms = append(forms, form{content: *widget.NewForm(
		widget.NewFormItem("1", widget.NewEntry()),
		widget.NewFormItem("2", widget.NewEntry()),
	)})

	forms = append(forms, form{content: *widget.NewForm(
		widget.NewFormItem("3", widget.NewEntry()),
		widget.NewFormItem("4", widget.NewEntry()),
	)})

	formsContainer := container.NewGridWithColumns(3)

	for _, element := range forms {
		formsContainer.Add(layout.NewSpacer())
		formsContainer.Add(&element.content)
		fmt.Println(element.content.Items[0].Text)
	}

	container_compose := formsContainer
	// forms := []widget.NewForm(widget.NewFormItem("Text", widget.NewEntry()))
	// fmt.Println(forms.Items[0].Text)
	// forms := &widget.Form{W
	// 	Items: []*widget.FormItem{
	// 		{Text: "", Widget: widget.NewEntry()},
	// 		{Text: "", Widget: widget.NewEntry()},
	// 		{Text: "", Widget: widget.NewEntry()},
	// 	},
	// }

	// form := widget.NewForm(widget.NewFormItem(""))

	// container_compose := container.NewVBox(
	// 	container.NewHBox(
	// 		widget.NewLabel("Number of services/ containers:"),
	// 		nrOfServicesEntry,
	// 		widget.NewButton("Create form", func() {
	// 			if nrOfServicesEntry.Text != "" {

	// 				nrOfServices, err := strconv.Atoi(nrOfServicesEntry.Text)

	// 				if err != nil {
	// 					log.Fatal(err)
	// 				}
	// 				for i := 1; i <= nrOfServices; i++ {
	// 					// forms.Add(container.NewVBox(
	// 					// 	container.NewHBox(
	// 					// 		widget.NewLabel("Name of the service:"),
	// 					// 		widget.NewEntry(),
	// 					// 	),
	// 					// 	container.NewHBox(
	// 					// 		widget.NewLabel("Standard image or custom dockerfile:"),
	// 					// 		widget.NewEntry(),
	// 					// 	),
	// 					// 	container.NewHBox(
	// 					// 		widget.NewLabel("Name of the service"),
	// 					// 		widget.NewEntry(),
	// 					// 	),
	// 					// 	container.NewHBox(
	// 					// 		widget.NewCheck("Bind ports?", func(b bool) {}),
	// 					// 	),
	// 					// ))
	// 					forms.Add(widget.NewForm(
	// 						widget.NewFormItem(
	// 							"Name of the service:",
	// 							widget.NewEntry(),
	// 						),
	// 						widget.NewFormItem(
	// 							"Standard image or custom dockerfile:",
	// 							widget.NewRadioGroup([]string{"Image", "Custom"}, func(s string) {}),
	// 						),
	// 						widget.NewFormItem(
	// 							"Name of image/ path to dockerfile folder:",
	// 							widget.NewEntry(),
	// 						),
	// 						widget.NewFormItem(
	// 							"How many bound ports do you want?",
	// 							widget.NewCheck("", func(b bool) {}),
	// 						),
	// 					))
	// 				}
	// 			}
	// 		}),
	// 	),
	// 	// layout.NewSpacer(),
	// 	forms,
	// 	// layout.NewSpacer(),
	// 	// layout.NewSpacer(),
	// 	// layout.NewSpacer(),
	// 	widget.NewButton("Generate", func() {}),
	// 	layout.NewSpacer(),
	// )

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
					// TODO add theme here aswell
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
