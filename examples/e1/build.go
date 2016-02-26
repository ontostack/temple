package main

import (
	"io"

	. "github.com/ontostack/buildr"
)

//go:generate ../../temple gen/e1gen.go
func main() {
	builder := File("e1.exe")

	e1go := File("e1/e1.go").Depends(builder).Make(func(...TargetI) bool {
		if !Exists("e1") && !Mkdir("e1") {
			return false
		}
		if !FillFile("e1/e1.go", func(w io.Writer) bool { makeLoop(w); return true }) {
			return false
		}
		return InDir("e1", GoFmt)
	})

	e1exe := File("e1/e1.exe").Depends(e1go).Make(func(...TargetI) bool {
		return InDir("e1", GoBuild)
	})

	e1exe.Build()
}
