package main

import (
	"github.com/codeskyblue/go-sh"
	. "github.com/ontostack/buildr"
)

//go:generate ..\..\temple.exe gen\e1gen.go
func main() {
	builder := File("e1.exe")

	e1go := File("e1/e1.go").Depends(builder).Make(func(...TargetI) bool {
		if !Exists("e1") && !Mkdir("e1") {
			return false
		}
		if f, ok := CreateIfNotExists("e1/e1.go"); !ok {
			return false
		} else {
			func() {
				defer f.Close()
				makeLoop(f)
			}()
		}
		return InDir("e1", func() bool {
			return Cmd(sh.Command("go", "fmt"))
		})
	})

	e1exe := File("e1/e1.exe").Depends(e1go).Make(func(...TargetI) bool {
		return InDir("e1", func() bool {
			return Cmd(sh.Command("go", "build"))
		})
	})

	e1exe.Build()
}
