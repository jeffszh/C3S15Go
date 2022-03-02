package main

import (
	"embed"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"image/png"
	"strings"
	"time"
)

//go:embed images
var imgs embed.FS

func main() {
	var inTE, outTE *walk.TextEdit

	screenWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	screenHeight := win.GetSystemMetrics(win.SM_CYSCREEN)
	wndWidth := 800
	wndHeight := 600
	var mainWndP *walk.MainWindow
	mainWnd := MainWindow{
		Title:   "三炮十五兵 Go语言版",
		MinSize: Size{Width: 600, Height: 400},
		Bounds: Rectangle{
			X:      (int(screenWidth) - wndWidth) / 2,
			Y:      (int(screenHeight) - wndHeight) / 2,
			Width:  wndWidth,
			Height: wndHeight,
		},
		Font:   Font{Family: "宋体", PointSize: 20},
		Layout: VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "SCREAM",
				OnClicked: func() {
					_ = outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
		OnMouseDown: func(x, y int, button walk.MouseButton) {
		},
		OnSizeChanged: func() {
			fmt.Printf("X=%d，Y=%d，宽度=%d，高度=%d\n",
				mainWndP.X(), mainWndP.Y(),
				mainWndP.Width(), mainWndP.Height())
		},
		OnBoundsChanged: func() {
			fmt.Println("OnBoundsChanged: " +
				fmt.Sprint(mainWndP.Bounds()))
		},
		AssignTo: &mainWndP,
	}
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println(mainWndP.Bounds())
		mainWnd.Bounds.X = 300
		fmt.Println(mainWndP.Bounds())
	}()
	pngFile, _ := imgs.Open("images/block.png")
	img, _ := png.Decode(pngFile)
	icon, _ := walk.NewIconFromImageForDPI(img, 300)
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = mainWndP.SetIcon(icon)
	}()
	_, _ = mainWnd.Run()
}
