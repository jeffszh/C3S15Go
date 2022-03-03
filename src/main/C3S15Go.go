package main

import (
	"embed"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

//go:embed images
var imgs embed.FS

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
