package main

import (
	"embed"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"image/jpeg"
	"image/png"
	"main/model"
	"time"
)

//go:embed images
var imgs embed.FS

//os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\SIMYOU.TTF")

func main() {
	bkImg, _ := jpeg.Decode(noErr(imgs.Open("images/wood.jpg")))
	bkBmp, _ := walk.NewBitmapFromImageForDPI(bkImg, 96)

	screenWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	screenHeight := win.GetSystemMetrics(win.SM_CYSCREEN)
	wndWidth := 800
	wndHeight := 600
	var mainWndPtr *walk.MainWindow
	cbs := NewChessBoard(&mainWndPtr, bkBmp)
	mainWnd := MainWindow{
		Title:   model.AppConfig.AppTitle,
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
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "按钮，很多个字的按钮",
					},
					HSpacer{},
				},
			},
			cbs.Declare(),
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "SCREAM",
					},
					HSpacer{},
				},
			},
		},
		OnMouseDown: func(x, y int, button walk.MouseButton) {
		},
		OnSizeChanged: func() {
			fmt.Printf("X=%d，Y=%d，宽度=%d，高度=%d\n",
				cbs.MainWnd().X(), cbs.MainWnd().Y(),
				cbs.MainWnd().Width(), cbs.MainWnd().Height())
		},
		OnBoundsChanged: func() {
			fmt.Println("OnBoundsChanged: " +
				fmt.Sprint(cbs.MainWnd().Bounds()))
		},
		AssignTo: &mainWndPtr,
	}
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println(cbs.MainWnd().Bounds())
		mainWnd.Bounds.X = 300
		fmt.Println(cbs.MainWnd().Bounds())
	}()
	pngFile, _ := imgs.Open("images/block.png")
	img, _ := png.Decode(pngFile)
	icon, _ := walk.NewIconFromImageForDPI(img, 300)
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = cbs.MainWnd().SetIcon(icon)
	}()
	_, _ = mainWnd.Run()
}
