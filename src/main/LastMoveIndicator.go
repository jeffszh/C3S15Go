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
)

type lastMoveIndicator struct {
	canvas.Image
	direction Direction
}

// LastMoveIndicator 表示最后一步的箭头。
//
// 这里再次看到Go语言不是面向对象所带来的麻烦，
// 更要命的是fyne竟然没有多边形，
// 而且搞不清楚fyne画图内部是怎样实现的（可能跟平台相关），
// 于是也就无法自己实现多边形了，只好用图片来做了。
type LastMoveIndicator interface {
	GetImage() *canvas.Image
	Direction() Direction
	SetDirection(direction Direction)
}

func NewLastMoveIndicator() LastMoveIndicator {
	println("==========================================")
	lmi := &lastMoveIndicator{
		Image: *canvas.NewImageFromResource(leftArrowPictureResource),
	}
	//lmi.Translucency = 0.5
	return lmi
}

func (l *lastMoveIndicator) GetImage() *canvas.Image {
	return &l.Image
}

func (l *lastMoveIndicator) Direction() Direction {
	return l.direction
}

func (l *lastMoveIndicator) SetDirection(direction Direction) {
	l.direction = direction
	// TODO 根据方向选择不同的图片
}
