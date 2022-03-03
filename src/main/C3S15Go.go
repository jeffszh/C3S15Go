package main

import (
	"embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"main/theme"
)

//go:embed images
var imgs embed.FS

func main() {
	//os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\SIMYOU.TTF")
	fmt.Println(theme.CurrentTheme)
	a := app.New()
	w := a.NewWindow("Hello 中文 World")

	w.SetContent(widget.NewLabel("Hello World! 这是汉字。"))
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
