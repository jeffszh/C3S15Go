package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
	"time"
)

func main() {
	var inTE, outTE *walk.TextEdit

	mainWnd := MainWindow{
		Title:   "输入几个中文吧。",
		MinSize: Size{600, 400},
		Layout:  VBox{},
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
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
		OnMouseDown: func(x, y int, button walk.MouseButton) {
		},
	}
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println(mainWnd.Bounds)
		mainWnd.Bounds.X = 300
		fmt.Println(mainWnd.Bounds)
	}()
	mainWnd.Run()
}
