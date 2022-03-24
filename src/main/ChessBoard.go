package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"image/png"
	"main/model"
	"time"
)

const (
	chessDiameterRatio   = 0.8
	chessTextSizeRatio   = 0.44
	chessTextOffsetRatio = 0.098

	chessBoardMinWidth = 300
	borderWidth        = 20
)

var (
	leftArrowImg,
	rightArrowImg,
	upArrowImg,
	downArrowImg walk.Image
)

func init() {
	leftArrowImg = loadImage("images/left.png")
	rightArrowImg = loadImage("images/right.png")
	upArrowImg = loadImage("images/up.png")
	downArrowImg = loadImage("images/down.png")
}

func loadImage(resourcePath string) walk.Image {
	file, _ := imgs.Open(resourcePath)
	img, _ := png.Decode(file)
	bmp, _ := walk.NewBitmapFromImageForDPI(img, 96)
	return bmp
}

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

	// 棋局状态改变时向上通知的回调
	onStatusChange func()
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
func NewChessBoard(mainWndToBeAssign **walk.MainWindow, chessBoardBackground *walk.Bitmap,
	onStatusChange func()) ChessBoard {
	cbs := chessBoardStruct{
		mainWnd:              mainWndToBeAssign,
		chessBoardBackground: chessBoardBackground,
		onStatusChange:       onStatusChange,
	}
	cbs.AssignTo = &cbs.chessBoard
	cbs.ClearsBackground = false
	cbs.InvalidatesOnResize = true
	cbs.MinSize = declarative.Size{Width: chessBoardMinWidth, Height: chessBoardMinWidth}
	cbs.Paint = cbs.OnPaint
	cbs.scene = model.NewScene()
	cbs.scene.SetOnChange(func(scene model.Scene) {
		_ = cbs.chessBoard.Invalidate()
		onStatusChange()
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
			if chess.Visible() &&
				// 并且不是正在拖动中的棋子
				!(draggingChess != nil && startDragCellX == cellX && startDragCellY == cellY) {
				drawChess(canvas,
					orgX+cellX*cellSize+cellSize/2,
					orgY+cellY*cellSize+cellSize/2,
					cellSize, chess)
			}
		}
	}

	// 画最后一步箭头
	//canvas.DrawPolylinePixels()
	// lxn/walk有画多边形，但没有填充多边形，所以仍然是要用图片来做。
	// 在这方面lxn/walk并不比fyne优胜。
	lastMove := cbs.Scene().LastMove()
	if lastMove != nil {
		fromX, fromY := lastMove.FromXY()
		toX, toY := lastMove.ToXY()
		var arrowImg walk.Image
		switch {
		case fromX > toX:
			arrowImg = leftArrowImg
		case fromX < toX:
			arrowImg = rightArrowImg
		case fromY > toY:
			arrowImg = upArrowImg
		case fromY < toY:
			arrowImg = downArrowImg
		}
		deltaX := toX - fromX
		deltaY := toY - fromY
		var xScale, yScale int
		if deltaX == 2 || deltaX == -2 {
			xScale = 2
		} else {
			xScale = 1
		}
		if deltaY == 2 || deltaY == -2 {
			yScale = 2
		} else {
			yScale = 1
		}
		rect := walk.Rectangle{
			X:      orgX,
			Y:      orgY,
			Width:  cellSize * xScale,
			Height: cellSize * yScale,
		}
		if fromX != toX {
			rect.X += cellSize / 2
		}
		if fromY != toY {
			rect.Y += cellSize / 2
		}
		minX := min(fromX, toX)
		minY := min(fromY, toY)
		rect.X += minX * cellSize
		rect.Y += minY * cellSize
		_ = canvas.DrawImageStretchedPixels(arrowImg, rect)
	}

	// 若正在拖动，画拖动影像。
	if draggingChess != nil {
		drawChess(canvas, mouseX, mouseY, cellSize, draggingChess)
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

/*
# 拖放的实现方法

由于lxn/walk提供的布局方式非常少而且死板（不要说自由定位了，就最简单的放置两个重叠的控件我也找不到办法），
即使克服了在运行时动态创建Widget的难题，也无法在运行时动态设置Widget的位置和大小，
注定必然只能通过OnPaint来产生拖放影像，这样，效率很低且会闪烁，效果极差。
当然，如果非要实现好的拖放效果，想想曲折些的办法还是做得到，
再说，也不是非要用拖放的方式来走棋，可以做点击的方式。
不过，现在不是非要将效果做好，而是通过做来体验这些东西是否好用，那么，lxn/walk大大失分了。
虽然fyne更加难用又难理解，不过fyne至少比lxn/walk考虑得周全些，而且也做得比较完整。
*/

// 拖放状态
var draggingChess model.Chess = nil
var startDragCellX, startDragCellY = -1, -1
var mouseX, mouseY int

func (cbs *chessBoardStruct) mouseDown(x, y int, button walk.MouseButton) {
	//fmt.Printf("mouse down: %d, %d\n", x, y)
	if button != walk.LeftButton {
		return
	}
	if cbs.Scene().GameOver() {
		return
	}
	mouseX = x
	mouseY = y
	cellX, cellY := cbs.mouseXyToCellXy(x, y)
	if model.AllInRange(cellX, cellY) {
		cellInd := model.XyToIndex(cellX, cellY)
		chess := cbs.Scene().ChessList()[cellInd]
		if chess.Type() == cbs.Scene().MovingSide() {
			draggingChess = chess
			startDragCellX = cellX
			startDragCellY = cellY
			_ = cbs.chessBoard.Invalidate()
		}
	}
}

func (cbs *chessBoardStruct) mouseMove(x, y int, _ walk.MouseButton) {
	//fmt.Printf("mouse move: %d, %d\n", x, y)
	if draggingChess != nil {
		mouseX = x
		mouseY = y
		_ = cbs.chessBoard.Invalidate()
	}
}

func (cbs *chessBoardStruct) mouseUp(x, y int, button walk.MouseButton) {
	//fmt.Printf("mouse up: %d, %d\n", x, y)
	if button != walk.LeftButton {
		return
	}
	if draggingChess == nil {
		return
	}

	// 走棋
	toX, toY := cbs.mouseXyToCellXy(x, y)
	move := model.NewMoveByXY(startDragCellX, startDragCellY, toX, toY)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cbs.Scene().ApplyMove(move)
	}()

	draggingChess = nil
	_ = cbs.chessBoard.Invalidate()
}

func (cbs *chessBoardStruct) mouseXyToCellXy(mouseX, mouseY int) (cellX, cellY int) {
	cellX = (mouseX - cbs.orgX) / cbs.cellSize
	cellY = (mouseY - cbs.orgY) / cbs.cellSize
	return
}
