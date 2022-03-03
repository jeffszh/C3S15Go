package main

import (
	"fmt"
	"io/fs"
)

type nn struct {
	f1 string `json:"f_1,omitempty" :"f_1"`
	f2 int    `json:"f_2"`
}

func noErr(p1 fs.File, _ error) fs.File {
	fmt.Println()
	return p1
}
