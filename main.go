package main

import (
	"flag"
	"github.com/ontostack/temple/temple"
)

func main() {
	// src is the input that we want to tokenize.
	flag.Parse()
	args := flag.Args()
	temple.Run(args...)
}
