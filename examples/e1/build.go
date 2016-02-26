package main

import (
	"os"

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
		if !Check(os.Chdir("e1")) {
			return false
		}
		if !Cmd(sh.Command("go", "fmt")) {
			return false
		}
		if !Check(os.Chdir("..")) {
			return false
		}
		return true
	})

	e1exe := File("e1/e1.exe").Depends(e1go).Make(func(...TargetI) bool {
		if !Check(os.Chdir("e1")) {
			return false
		}
		if !Cmd(sh.Command("go", "build")) {
			return false
		}
		if !Check(os.Chdir("..")) {
			return false
		}
		return true
	})

	e1exe.Build()
}
