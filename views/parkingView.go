package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func SetupWindow(mainApp fyne.App) (fyne.Window, *fyne.Container) {
	win := mainApp.NewWindow("Parking Lot Simulator")
	win.Resize(fyne.NewSize(800, 500))
	win.SetFixedSize(true)

	imgResource, _ := fyne.LoadResourceFromPath("assets/bg.png")
	img := canvas.NewImageFromResource(imgResource)
	img.Resize(fyne.NewSize(800, 500))

	cont := container.NewWithoutLayout(img)
	win.SetContent(cont)
	return win, cont
}
