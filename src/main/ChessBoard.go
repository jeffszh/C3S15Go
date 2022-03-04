package main

import (
	"embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

//go:embed images
var imgs embed.FS

type chessBoard struct {
	widget.BaseWidget
}

type chessBoardRenderer struct {
	chessBoard *chessBoard
}

func NewChessBoard() *chessBoard {
	return &chessBoard{}
}

func (*chessBoard) CreateRenderer() fyne.WidgetRenderer {
	return &chessBoardRenderer{}
}

func (c chessBoardRenderer) Destroy() {
	println("Destroy()")
}

func (c chessBoardRenderer) Layout(size fyne.Size) {
	println("布局")
}

func (c chessBoardRenderer) MinSize() fyne.Size {
	return 
}

func (c chessBoardRenderer) Objects() []fyne.CanvasObject {
	panic("implement me")
}

func (c chessBoardRenderer) Refresh() {
	panic("implement me")
}
