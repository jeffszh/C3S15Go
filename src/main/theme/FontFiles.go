package theme

import (
	"embed"
	"fyne.io/fyne"
)

//go:embed fonts/simsun.ttc
var myFonts embed.FS

var resourceSimsunTtc = &fyne.StaticResource{
	StaticName:    "simsun.ttc",
	StaticContent: loadFont("fonts/simsun.ttc"),// 这个可以。
	//StaticContent: loadFont("fonts/SIMYOU_7.ttf"),// 这个不行，出现“Fyne error:  font load error”
	//StaticContent: loadFont("fonts/msyhl.ttc"),// 这个也不行。
	// 看来fyne对字体的支持不太好，大多数字体都是不能正常加载的，
	// 但是用环境变量指定字体却可行，例如：os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\SIMYOU.TTF")
}

func loadFont(f string) []byte {
	font, _ := myFonts.ReadFile(f)
	return font
}
