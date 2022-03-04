// Package theme theme.go
package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"image/color"
)

type MyTheme struct{}

var AppTheme fyne.Theme = &MyTheme{}

// Font return bundled font resource
// ResourceSourceHanSansTtf 即是 bundle.go 文件中 var 的变量名
func (m MyTheme) Font(fyne.TextStyle) fyne.Resource {
	return resourceSimsunTtc
}

func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(name fyne.ThemeSizeName) float32 {
	//return theme.DefaultTheme().Size(n)
	//println(n, theme.DefaultTheme().Size(n))
	switch name {
	case "text":
		return 20
	default:
		return theme.DefaultTheme().Size(name)
	}
}
