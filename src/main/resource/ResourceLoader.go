package resource

import "embed"

//go:embed fonts/simsun.ttc images
var embFiles embed.FS

func LoadFileInBytes(filename string) []byte {
	bytes, _ := embFiles.ReadFile(filename)
	return bytes
}
