package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"main/resource"
)

var backgroundPictureResource = &fyne.StaticResource{
	StaticName:    "background",
	StaticContent: resource.LoadFileInBytes("images/wood.jpg"),
}

type chessBoard struct {
	widget.BaseWidget
	background fyne.CanvasObject
}

type chessBoardRenderer struct {
	chessBoard *chessBoard
}

func NewChessBoard() *chessBoard {
	//bg := canvas.NewRectangle(color.RGBA{R: 255, A: 255})
	bg := canvas.NewImageFromResource(backgroundPictureResource)
	cb := chessBoard{background: bg}
	return &cb
}

func (cb *chessBoard) CreateRenderer() fyne.WidgetRenderer {
	return &chessBoardRenderer{cb}
}

func (cbr *chessBoardRenderer) Destroy() {
	fmt.Println("Destroy()")
}

func (cbr *chessBoardRenderer) Layout(size fyne.Size) {
	fmt.Printf("布局，size = %f x %f\n", size.Width, size.Height)
	cbr.chessBoard.background.Resize(fyne.NewSize(300, 200))
	cbr.chessBoard.background.Move(fyne.NewPos(30, 20))
}

func (cbr *chessBoardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(80, 60)
}

func (cbr *chessBoardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{cbr.chessBoard.background}
}

func (cbr *chessBoardRenderer) Refresh() {
	fmt.Println("刷新。")
	cbr.chessBoard.Refresh()
	cbr.chessBoard.background.Refresh()
}
