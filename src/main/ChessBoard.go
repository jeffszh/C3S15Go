package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type chessBoardStruct struct {
	declarative.CustomWidget
	mainWnd              **walk.MainWindow
	chessBoard           *walk.CustomWidget
	chessBoardBackground *walk.Bitmap
}

type ChessBoard interface {
	Declare() *declarative.CustomWidget
	MainWnd() *walk.MainWindow
	OnPaint(canvas *walk.Canvas, updateBounds walk.Rectangle) error
}

// NewChessBoard 创建ChessBoard
// 参数：mainWndToBeAssign 即将赋值 mainWnd 的指针变量的指针
//（这是lxn/walk的奇特之处，declarative的时候并不知道，运行时才知道赋值到哪了。）
// 参数：chessBoardBackground 背景图片
func NewChessBoard(mainWndToBeAssign **walk.MainWindow, chessBoardBackground *walk.Bitmap) ChessBoard {
	cbs := chessBoardStruct{
		mainWnd:              mainWndToBeAssign,
		chessBoardBackground: chessBoardBackground,
	}
	cbs.AssignTo = &cbs.chessBoard
	cbs.ClearsBackground = false
	cbs.InvalidatesOnResize = true
	cbs.Paint = cbs.OnPaint
	return &cbs
}

// Declare 声明
// 返回声明期的结构体
func (cbs *chessBoardStruct) Declare() *declarative.CustomWidget {
	return &cbs.CustomWidget
}

// MainWnd 取运行期的主窗口
func (cbs *chessBoardStruct) MainWnd() *walk.MainWindow {
	return *cbs.mainWnd
}

func (cbs *chessBoardStruct) OnPaint(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	bounds := cbs.chessBoard.Bounds()
	bounds.X = 0
	bounds.Y = 0
	_ = canvas.DrawImageStretchedPixels(cbs.chessBoardBackground, bounds)
	return nil
}
