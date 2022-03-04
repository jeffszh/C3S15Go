package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"main/theme"
	"time"
)

func main() {
	//os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\SIMYOU.TTF")
	appMain := app.New()
	appMain.Settings().SetTheme(theme.AppTheme)
	mainWnd := appMain.NewWindow("三炮十五兵 Go语言版")

	statusText := widget.NewLabel("就绪")
	bottomPane := container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(), statusText, layout.NewSpacer())
	btnRestart := widget.NewButton("重新开始", func() {
		fmt.Println("按钮：重新开始")
		statusText.SetText("重新开始")
		go func() {
			time.Sleep(3 * time.Second)
			statusText.SetText("就绪")
		}()
	})
	topPane := container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(), btnRestart, layout.NewSpacer())
	layout.NewBorderLayout(nil, nil, nil, nil)
	mainWnd.SetContent(container.New(layout.NewBorderLayout(
		topPane, bottomPane, nil, nil),
		topPane, bottomPane))
	// 无法指定窗口位置，但是却有居中，真奇葩。
	mainWnd.CenterOnScreen()
	mainWnd.Resize(fyne.NewSize(1200, 900))
	go func() {
		// 通过对中和改变尺寸，实际上可以把窗口移到任意位置，有必要的时候就用这招吧。
		time.Sleep(100 * time.Millisecond)
		mainWnd.Resize(fyne.NewSize(800, 600))
	}()
	mainWnd.ShowAndRun()
}
