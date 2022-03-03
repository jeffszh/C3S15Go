package theme

import (
	"embed"
	"fyne.io/fyne"
)

//go:embed fonts
var myFonts embed.FS

var resourceSimsunTtc = &fyne.StaticResource{
	StaticName:    "simsun.ttc",
	StaticContent: loadFont("fonts/simsun.ttc"),
}

func loadFont(f string) []byte {
	font, _ := myFonts.ReadFile(f)
	return font
}
