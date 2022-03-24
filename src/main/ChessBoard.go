package main

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"main/model"
)

const (
	chessDiameterRatio   = 0.8
	chessTextSizeRatio   = 0.44
	chessTextOffsetRatio = 0.098

	chessBoardMinWidth = 300
	borderWidth        = 20
)

// 棋盘内部结构
type chessBoardStruct struct {
	declarative.CustomWidget
	mainWnd              **walk.MainWindow
	chessBoard           *walk.CustomWidget
	chessBoardBackground *walk.Bitmap

	// 棋盘中央正方形区域原点
	orgX, orgY int
	// 棋盘格大小
	cellSize int

	scene model.Scene
}

// ChessBoard 棋盘
type ChessBoard interface {
	Declare() *declarative.CustomWidget
	MainWnd() *walk.MainWindow
	OnPaint(canvas *walk.Canvas, updateBounds walk.Rectangle) error
	Scene() model.Scene
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
	cbs.MinSize = declarative.Size{Width: chessBoardMinWidth, Height: chessBoardMinWidth}
	cbs.Paint = cbs.OnPaint
	cbs.scene = model.NewScene()
	cbs.scene.SetOnChange(func(scene model.Scene) {
		_ = cbs.chessBoard.Invalidate()
	})
	cbs.OnMouseDown = cbs.mouseDown
	cbs.OnMouseMove = cbs.mouseMove
	cbs.OnMouseUp = cbs.mouseUp
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

/*
func (cbs *chessBoardStruct) OnPaint(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	printRect(updateBounds)
	printRect(cbs.chessBoard.ClientBounds())
	bounds := cbs.chessBoard.Bounds()
	bounds.X = 0
	bounds.Y = 0
	_ = canvas.DrawImageStretchedPixels(cbs.chessBoardBackground, bounds)
	return nil
}

func printRect(rectangle walk.Rectangle) {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Printf("%d,%d,%d,%d\n", rectangle.X, rectangle.Y, rectangle.Width, rectangle.Height)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>")
}
*/

// golang竟然没有两整数求最小值的内置函数。
func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func (cbs *chessBoardStruct) OnPaint(canvas *walk.Canvas, _ walk.Rectangle) error {
	bounds := cbs.chessBoard.ClientBounds()
	_ = canvas.DrawImageStretchedPixels(cbs.chessBoardBackground, bounds)

	// 计算尺寸
	minXy := min(bounds.Width, bounds.Height)
	cellSize := (minXy - 2*borderWidth) / 5
	orgX := (bounds.Width - cellSize*5) / 2
	orgY := (bounds.Height - cellSize*5) / 2

	// 存起来
	cbs.cellSize = cellSize
	cbs.orgX = orgX
	cbs.orgY = orgY

	// 横竖线
	brush, _ := walk.NewSolidColorBrush(0)
	thickerPen, _ := walk.NewGeometricPen(walk.PenSolid, 3, brush)
	thinnerPen, _ := walk.NewGeometricPen(walk.PenSolid, 2, brush)
	for i := 0; i <= 5; i++ {
		var pen walk.Pen
		if i == 0 || i == 5 {
			pen = thickerPen
		} else {
			pen = thinnerPen
		}
		// 横线
		_ = canvas.DrawLinePixels(pen,
			walk.Point{X: orgX, Y: orgY + i*cellSize},
			walk.Point{X: orgX + 5*cellSize, Y: orgY + i*cellSize})
		// 竖线
		_ = canvas.DrawLinePixels(pen,
			walk.Point{X: orgX + i*cellSize, Y: orgY},
			walk.Point{X: orgX + i*cellSize, Y: orgY + 5*cellSize})
	}
	defer thinnerPen.Dispose()
	defer thickerPen.Dispose()
	defer brush.Dispose()

	// 画棋子
	chessList := cbs.scene.ChessList()
	for cellY := 0; cellY < 5; cellY++ {
		for cellX := 0; cellX < 5; cellX++ {
			ind := model.XyToIndex(cellX, cellY)
			chess := chessList[ind]
			visible := chess.Visible()
			if startDragX == cellX && startDragY == cellY {
				visible = false
			}

			if visible {
				drawChess(canvas,
					orgX+cellX*cellSize+cellSize/2,
					orgY+cellY*cellSize+cellSize/2,
					cellSize, chess)
			}
		}
	}
	return nil
}

func (cbs *chessBoardStruct) Scene() model.Scene {
	return cbs.scene
}

func drawChess(canvas *walk.Canvas, centerX, centerY, cellSize int, chess model.Chess) {
	r, g, b, _ := chess.Color().RGBA()
	cellColor := walk.RGB(byte(r), byte(g), byte(b))
	brush, _ := walk.NewSolidColorBrush(cellColor)
	pen, _ := walk.NewGeometricPen(walk.PenSolid, 3, brush)
	whiteBrush, _ := walk.NewSolidColorBrush(walk.RGB(255, 255, 255))
	text := chess.Text()

	widthAndHeight := int(float64(cellSize) * chessDiameterRatio)
	x := centerX - widthAndHeight/2
	y := centerY - widthAndHeight/2

	// 圆圈
	rect := walk.Rectangle{
		X:      x,
		Y:      y,
		Width:  widthAndHeight,
		Height: widthAndHeight,
	}
	_ = canvas.FillEllipsePixels(whiteBrush, rect)
	_ = canvas.DrawEllipsePixels(pen, rect)

	// 文字
	font, _ := walk.NewFont("楷体", int(float64(cellSize)*chessTextSizeRatio),
		//font, _ := walk.NewFont("宋体", int(float64(cellSize)*chessTextSizeRatio),
		walk.FontBold)
	rect.Y += int(float64(cellSize) * chessTextOffsetRatio)
	_ = canvas.DrawTextPixels(text, font, cellColor, rect, walk.TextCenter)
	font.Dispose()

	pen.Dispose()
	brush.Dispose()
	whiteBrush.Dispose()
}

// 拖放状态
var dragging = false
var startDragX, startDragY = -1, -1

func (cbs *chessBoardStruct) mouseDown(x, y int, button walk.MouseButton) {
	fmt.Printf("mouse down: %d, %d\n", x, y)
}

func (cbs *chessBoardStruct) mouseMove(x, y int, button walk.MouseButton) {
	fmt.Printf("mouse move: %d, %d\n", x, y)
}

func (cbs *chessBoardStruct) mouseUp(x, y int, button walk.MouseButton) {
	fmt.Printf("mouse up: %d, %d\n", x, y)
}
