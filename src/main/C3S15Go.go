package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"main/model"
	"main/theme"
)

var mainChessBoard = NewChessBoard()
var statusText *widget.Label

func main() {
	//os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\SIMYOU.TTF")
	appMain := app.New()
	appMain.Settings().SetTheme(theme.AppTheme)
	mainWnd := appMain.NewWindow(model.AppConfig.AppTitle)

	statusText = widget.NewLabel("就绪")
	bottomPane := container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(), statusText, layout.NewSpacer())
	btnRestart := widget.NewButton("重新开始", func() {
		//dialog.NewInformation("提示消息", "马上重新开始。", mainWnd).Show()
		//fmt.Println("按钮：重新开始")
		//statusText.SetText("重新开始")
		//go func() {
		//	time.Sleep(3 * time.Second)
		//	statusText.SetText("就绪")
		//}()
		restartGame()
	})
	topPane := container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(), btnRestart, layout.NewSpacer())
	layout.NewBorderLayout(nil, nil, nil, nil)
	mainWnd.SetContent(container.New(layout.NewBorderLayout(
		topPane, bottomPane, nil, nil),
		topPane, bottomPane, mainChessBoard))
	mainChessBoard.onGameInfoChanged = updateStatusText
	mainWnd.Resize(fyne.NewSize(800, 600))
	// 无法指定窗口位置，但是却有居中，真奇葩。
	mainWnd.CenterOnScreen()
	restartGame()
	mainWnd.ShowAndRun()
}

func restartGame() {
	mainChessBoard.scene.SetInitialContent()
}

func updateStatusText() {
	conf := model.AppConfig
	mainScene := mainChessBoard.scene
	statusText.Text = fmt.Sprintf(
		"%s：%s  %s：%s    步数：%d  轮到【%s】走棋",
		conf.SoldierText, model.PlayerTypeText(conf.SoldierPlayType),
		conf.CannonText, model.PlayerTypeText(conf.SoldierPlayType),
		mainScene.MoveCount(), model.NewChess(mainScene.MovingSide()).Text(),
	)
	statusText.Refresh()
}
