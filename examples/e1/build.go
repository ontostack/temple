package main

import (
	"os"
	"log"
)

//go:generate ..\..\temple.exe gen\e1.go
func main() {
	f, err := os.Create(os.Args[1])
	if os.IsExist(err) {
		f, err = os.Open(os.Args[1])
	}
	if err != nil {
		log.Fatalln(err)
	}
	makeLoop(f)
}
