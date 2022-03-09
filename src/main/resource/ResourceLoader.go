package resource

import "embed"

//go:embed fonts/simsun.ttc images
var embFiles embed.FS

func LoadFileInBytes(f string) []byte {
	font, _ := embFiles.ReadFile(f)
	return font
}
