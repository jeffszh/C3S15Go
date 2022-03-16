package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"main/resource"
)

type Direction int

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionUp
	DirectionDown
)

var (
	leftArrowPictureResource = &fyne.StaticResource{
		StaticName:    "leftArrow",
		StaticContent: resource.LoadFileInBytes("images/left.png"),
	}
	rightArrowPictureResource = &fyne.StaticResource{
		StaticName:    "rightArrow",
		StaticContent: resource.LoadFileInBytes("images/right.png"),
	}
	upArrowPictureResource = &fyne.StaticResource{
		StaticName:    "upArrow",
		StaticContent: resource.LoadFileInBytes("images/up.png"),
	}
	downArrowPictureResource = &fyne.StaticResource{
		StaticName:    "downArrow",
		StaticContent: resource.LoadFileInBytes("images/down.png"),
	}

	leftArrow  = canvas.NewImageFromResource(leftArrowPictureResource)
	rightArrow = canvas.NewImageFromResource(rightArrowPictureResource)
	upArrow    = canvas.NewImageFromResource(upArrowPictureResource)
	downArrow  = canvas.NewImageFromResource(downArrowPictureResource)
)

type lastMoveIndicator struct {
	image     *canvas.Image
	direction Direction
}

// LastMoveIndicator 表示最后一步的箭头。
//
// 这里再次看到Go语言不是面向对象所带来的麻烦，
// 更要命的是fyne竟然没有多边形，
// 而且搞不清楚fyne画图内部是怎样实现的（可能跟平台相关），
// 于是也就无法自己实现多边形了，只好用图片来做了。
type LastMoveIndicator interface {
	Image() *canvas.Image
	Direction() Direction
	SetDirection(direction Direction)
}

func NewLastMoveIndicator() LastMoveIndicator {
	println("==========================================")
	lmi := &lastMoveIndicator{
		image: leftArrow,
	}
	//lmi.Translucency = 0.5
	return lmi
}

func (l *lastMoveIndicator) Image() *canvas.Image {
	return l.image
}

func (l *lastMoveIndicator) Direction() Direction {
	return l.direction
}

func (l *lastMoveIndicator) SetDirection(direction Direction) {
	l.direction = direction
	switch direction {
	case DirectionLeft:
		l.image = leftArrow
	case DirectionRight:
		l.image = rightArrow
	case DirectionUp:
		l.image = upArrow
	case DirectionDown:
		l.image = downArrow
	}
}

func repositionMoveIndicator(lmi LastMoveIndicator, board *chessBoard) {
	scene := board.scene
	lastMove := scene.LastMove()
	if lastMove == nil {
		lmi.Image().Hide()
		return
	}
	lmi.Image().Show()

	fromX, fromY := lastMove.FromXY()
	toX, toY := lastMove.ToXY()
	var dir Direction
	switch {
	case fromX > toX:
		dir = DirectionLeft
	case fromX < toX:
		dir = DirectionRight
	case fromY > toY:
		dir = DirectionUp
	case fromY < toY:
		dir = DirectionDown
	}
	lmi.SetDirection(dir)

	deltaX := toX - fromX
	deltaY := toY - fromY
	var xScale, yScale float32
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
	lmi.Image().Resize(fyne.NewSize(board.cellSize*xScale, board.cellSize*yScale))

	startX := (board.Size().Width - board.cellSize*5) / 2
	startY := (board.Size().Height - board.cellSize*5) / 2
	if fromX != toX {
		startX += board.cellSize / 2
	}
	if fromY != toY {
		startY += board.cellSize / 2
	}
	var minX, minY float32
	if fromX < toX {
		minX = float32(fromX)
	} else {
		minX = float32(toX)
	}
	if fromY < toY {
		minY = float32(fromY)
	} else {
		minY = float32(toY)
	}
	lmi.Image().Move(fyne.NewPos(
		startX+minX*board.cellSize,
		startY+minY*board.cellSize,
	))
}
