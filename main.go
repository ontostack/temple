package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ontostack/temple/temple"
)

func main() {
	// src is the input that we want to tokenize.
	flag.Parse()
	args := flag.Args()
	log.Println("os.Args:", os.Args)
	log.Println("flag.Args:", args)
	log := log.New(os.Stderr, "", log.LstdFlags)
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}
	srcName := args[0]
	dest, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln(err)
	}
	if len(args) > 1 {
		dest, err = filepath.Abs(args[1])
		if err != nil {
			log.Fatalln(err)
		}
	}
	stat, err := os.Stat(dest)
	if os.IsNotExist(err) {
		log.Fatalln(err)
	}
	if !stat.IsDir() {
		log.Fatalln("destination should be directory path")
	}
	src, err := ioutil.ReadFile(srcName)
	if err != nil {
		log.Fatalln(err)
	}
	outputName := filepath.Join(dest, filepath.Base(srcName))
	buf := &bytes.Buffer{}
	temple.New("sample.go", src, buf).Run()
	b, err := format.Source(buf.Bytes())
	if err != nil {
		ioutil.WriteFile(outputName, buf.Bytes(), os.ModePerm)
		log.Fatalln(err)
	}
	ioutil.WriteFile(outputName, b, os.ModePerm)
}
