package main

import (
	"github.com/lxn/walk"
)

func (mws *mainWndStuff) chessBoardOnPaint(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	bounds := mws.chessBoard.Bounds()
	bounds.X = 0
	bounds.Y = 0
	_ = canvas.DrawImageStretchedPixels(mws.chessBoardBackground, bounds)
	return nil
}
