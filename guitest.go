package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.Resize(fyne.NewSize(500, 600))
	label1 := widget.NewLabel("You are in content 1")
	var content1 fyne.CanvasObject
	var content2 fyne.CanvasObject
	button1 := widget.NewButton("Button 1",
		func() {
			fmt.Println("tapped button1")
			myWindow.SetContent(content2)
			myWindow.Show()
		})

	label2 := widget.NewLabel("You are in content 2")
	button2 := widget.NewButton("Button 2",
		func() {
			fmt.Println("tapped button2")
			myWindow.SetContent(content1)
			myWindow.Show()
		})
	content1 = container.NewVBox(
		label1,
		button1,
	)
	content2 = container.NewVBox(
		label2,
		button2,
	)
	myWindow.SetContent(content1)
	myWindow.ShowAndRun()
}
