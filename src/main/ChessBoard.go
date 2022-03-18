package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"main/model"
	"main/resource"
)

const (
	chessDiameterRatio = 0.8
	chessTextSizeRatio = 0.44

	chessBoardMinWidth = 200.0
)

var backgroundPictureResource = &fyne.StaticResource{
	StaticName:    "background",
	StaticContent: resource.LoadFileInBytes("images/wood.jpg"),
}

type chessBoard struct {
	widget.BaseWidget
	background fyne.CanvasObject
	//hLines     [7]*canvas.Line
	hLines [6]fyne.CanvasObject
	vLines [6]fyne.CanvasObject

	chessCells        [25]*chessCell
	lastMoveIndicator LastMoveIndicator

	cellSize   float32
	orgX, orgY float32

	scene model.Scene

	// 拖放操作
	draggingStarted    bool
	startDragCellIndex int
	draggingCell       *chessCell

	// 通知外面
	onGameInfoChanged func()
}

type chessBoardRenderer struct {
	chessBoard *chessBoard
}

type chessCell struct {
	circle *canvas.Circle
	text   *canvas.Text
}

func NewChessBoard() *chessBoard {
	//bg := canvas.NewRectangle(color.RGBA{R: 255, A: 255})
	bg := canvas.NewImageFromResource(backgroundPictureResource)
	cb := chessBoard{background: bg}
	for i := range cb.hLines {
		hLine := canvas.NewLine(color.Black)
		vLine := canvas.NewLine(color.Black)
		if i == 0 || i == 5 {
			hLine.StrokeWidth = 3
			vLine.StrokeWidth = 3
		} else {
			hLine.StrokeWidth = 2
			vLine.StrokeWidth = 2
		}
		cb.hLines[i] = hLine
		cb.vLines[i] = vLine
	}
	for i := range cb.chessCells {
		cell := chessCell{
			canvas.NewCircle(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
			canvas.NewText("炮", color.Black),
		}
		cell.circle.StrokeWidth = 3
		cell.circle.StrokeColor = color.Black
		//cell.text.Text = "兵"
		cell.text.Alignment = fyne.TextAlignCenter
		cb.chessCells[i] = &cell
	}
	cb.lastMoveIndicator = NewLastMoveIndicator()
	cb.lastMoveIndicator.Image().Resize(fyne.NewSize(120, 120))
	cb.lastMoveIndicator.Image().Move(fyne.NewPos(30, 30))

	// 初始化用于拖放的cell。
	cb.draggingStarted = false
	cb.startDragCellIndex = -1
	cb.draggingCell = &chessCell{
		canvas.NewCircle(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		canvas.NewText("炮", color.Black),
	}
	cb.draggingCell.circle.StrokeWidth = 3
	cb.draggingCell.circle.StrokeColor = color.Black
	cb.draggingCell.text.Alignment = fyne.TextAlignCenter
	cb.draggingCell.circle.Hide()
	cb.draggingCell.text.Hide()

	//cb.onGameInfoChanged = func() {
	//	// do nothing
	//}
	cb.scene = model.NewScene()
	cb.scene.SetOnChange(cb.applyScene)
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
	cbr.chessBoard.sizeChanged(size)
}

func (cb *chessBoard) sizeChanged(size fyne.Size) {
	// 防止出錯。
	if size.Width < chessBoardMinWidth || size.Height < chessBoardMinWidth {
		return
	}

	// 计算尺寸
	borderWidth := float32(20)
	cellSize := float32(int((fyne.Min(size.Width, size.Height) - 2*borderWidth) / 5))
	startX := (size.Width - cellSize*5) / 2
	startY := (size.Height - cellSize*5) / 2

	// 存入到chessBoard中。
	cb.cellSize = cellSize
	cb.orgX = startX
	cb.orgY = startY

	// 背景
	cb.background.Resize(fyne.NewSize(
		cellSize*5+borderWidth*2, cellSize*5+borderWidth*2))
	cb.background.Move(fyne.NewPos(
		startX-borderWidth, startY-borderWidth))

	// 横竖线
	for i := range cb.hLines {
		hLine := cb.hLines[i].(*canvas.Line) // 可以用类型断言
		hLine.Resize(fyne.NewSize(cellSize*5, 0))
		hLine.Move(fyne.NewPos(startX, startY+cellSize*float32(i)))
		vLine := cb.vLines[i].(*canvas.Line)
		vLine.Resize(fyne.NewSize(0, cellSize*5))
		vLine.Move(fyne.NewPos(startX+cellSize*float32(i), startY))
	}

	// 棋子
	for i := range cb.chessCells {
		x := i % 5
		y := i / 5
		cell := cb.chessCells[i]
		cell.circle.Resize(fyne.NewSize(cellSize*chessDiameterRatio, cellSize*chessDiameterRatio))
		cell.circle.Move(fyne.NewPos(
			startX+float32(x)*cellSize+cellSize*(1-chessDiameterRatio)/2,
			startY+float32(y)*cellSize+cellSize*(1-chessDiameterRatio)/2,
		))
		text := cell.text
		text.TextSize = cellSize * chessTextSizeRatio
		text.Resize(fyne.NewSize(cellSize, cellSize))
		text.Move(fyne.NewPos(startX+float32(x)*cellSize, startY+float32(y)*cellSize))
	}

	repositionMoveIndicator(cb.lastMoveIndicator, cb)
}

func (cb *chessBoard) Resize(size fyne.Size) {
	//fmt.Printf("resize: %f, %f\n", size.Width, size.Height)
	cb.sizeChanged(size)
	cb.BaseWidget.Resize(size)
}

func (cb *chessBoard) applyScene(scene model.Scene) {
	for i := range scene.ChessList() {
		ch := scene.ChessList()[i]
		chc := (cb.chessCells)[i]
		chc.text.Text = ch.Text()
		chc.text.Color = ch.Color()
		chc.text.Refresh()
		chc.text.Hidden = !ch.Visible()
		chc.circle.StrokeColor = ch.Color()
		chc.circle.Refresh()
		chc.circle.Hidden = !ch.Visible()
	}
	//cb.Refresh()
	repositionMoveIndicator(cb.lastMoveIndicator, cb)
	cb.onGameInfoChanged()
}

func (cbr *chessBoardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 150)
}

func (cbr *chessBoardRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{cbr.chessBoard.background}
	objs = append(objs, cbr.chessBoard.hLines[:]...)
	// 这里可以看到Go语言不面向对象，也没有泛型支持，总是会有缺陷的。
	// 如果 hLines 的类型为 *canvas.Line 的数组，那么只能用下面的循环语句来添加进来，
	// 如果类型改为 fyne.CanvasObject 的数组，可以整体引用了，但不知道具体类型。
	// 不同类型的数组或切片要转换，并没有好方法，只能不规矩地通过内部指针来做。
	//for i := range cbr.chessBoard.hLines {
	//	objs = append(objs, cbr.chessBoard.hLines[i])
	//}
	objs = append(objs, cbr.chessBoard.vLines[:]...)
	for i := range cbr.chessBoard.chessCells {
		objs = append(objs, cbr.chessBoard.chessCells[i].circle)
		objs = append(objs, cbr.chessBoard.chessCells[i].text)
	}
	objs = append(objs, cbr.chessBoard.lastMoveIndicator.Image())
	objs = append(objs, cbr.chessBoard.draggingCell.circle)
	objs = append(objs, cbr.chessBoard.draggingCell.text)
	return objs
}

func (cbr *chessBoardRenderer) Refresh() {
	fmt.Println("刷新。")
	//cbr.chessBoard.Refresh()
	//cbr.chessBoard.background.Refresh()
	//widget.BaseWidget.Refresh(cbr)
}

func (cb *chessBoard) Tapped(e *fyne.PointEvent) {
	fmt.Printf("点击！%f, %f\n", e.Position.X, e.Position.Y)
}

func (cb *chessBoard) Dragged(event *fyne.DragEvent) {
	mouseX, mouseY := event.Position.Components()
	dx, dy := event.Dragged.Components()
	cellX, cellY := cb.mouseXyToCellXy(mouseX, mouseY)
	fmt.Printf("拖动！%f,%f --> %f,%f | %d,%d\n", mouseX, mouseY, dx, dy, cellX, cellY)

	if !cb.draggingStarted {
		// 如果之前未drag，现在开始drag，不管位置是否合法，总要设置已开始drag的标志。
		cb.draggingStarted = true
		if !cb.scene.GameOver() &&
			0 <= cellX && cellX < 5 &&
			0 <= cellY && cellY < 5 {
			cellInd := model.XyToIndex(cellX, cellY)
			if cb.scene.ChessList()[cellInd].Type() == cb.scene.MovingSide() {
				// 记下开始拖动的位置
				cb.startDragCellIndex = cellInd
				// 原先位置上的cell隐藏掉
				cb.chessCells[cellInd].text.Hide()
				cb.chessCells[cellInd].circle.Hide()
				// 正在拖动的cell改成跟原先位置上的一模一样
				cb.draggingCell.text.Text = cb.chessCells[cellInd].text.Text
				cb.draggingCell.text.TextSize = cb.chessCells[cellInd].text.TextSize
				cb.draggingCell.text.Color = cb.chessCells[cellInd].text.Color
				cb.draggingCell.text.Resize(cb.chessCells[cellInd].text.Size())
				cb.draggingCell.circle.StrokeColor = cb.chessCells[cellInd].circle.StrokeColor
				cb.draggingCell.circle.Resize(cb.chessCells[cellInd].circle.Size())
				cb.draggingCell.text.Show()
				cb.draggingCell.circle.Show()
			}
		}
	}

	// 如果已经有合法的drag，更新位置。
	if cb.startDragCellIndex >= 0 {
		cb.draggingCell.circle.Move(fyne.NewPos(
			mouseX-cb.draggingCell.circle.Size().Width/2,
			mouseY-cb.draggingCell.circle.Size().Height/2,
		))
		cb.draggingCell.text.Move(fyne.NewPos(
			mouseX-cb.draggingCell.text.Size().Width/2,
			mouseY-cb.draggingCell.text.Size().Height/2,
		))
		cb.draggingCell.circle.Refresh()
		cb.draggingCell.text.Refresh()
	}
}

func (cb *chessBoard) DragEnd() {
	//println("拖完。")
	if cb.startDragCellIndex >= 0 {
		cb.chessCells[cb.startDragCellIndex].text.Show()
		cb.chessCells[cb.startDragCellIndex].circle.Show()
		cb.draggingCell.text.Hide()
		cb.draggingCell.circle.Hide()
	}

	// 走棋
	if cb.startDragCellIndex >= 0 {
		mx, my := cb.draggingCell.circle.Position().Components()
		sx, sy := cb.draggingCell.circle.Size().Components()
		mouseX := mx + sx/2
		mouseY := my + sy/2
		toX, toY := cb.mouseXyToCellXy(mouseX, mouseY)
		fromX, fromY := model.IndexToXY(cb.startDragCellIndex)
		if model.AllInRange(fromX, fromY, toX, toY) {
			move := model.NewMoveByXY(fromX, fromY, toX, toY)
			cb.scene.ApplyMove(move)
		}
	}

	cb.startDragCellIndex = -1
	cb.draggingStarted = false
}

func (cb *chessBoard) mouseXyToCellXy(mouseX, mouseY float32) (cellX, cellY int) {
	cellX = int((mouseX - cb.orgX) / cb.cellSize)
	cellY = int((mouseY - cb.orgY) / cb.cellSize)
	return
}
