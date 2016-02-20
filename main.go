package main

import (
	"bytes"
	"fmt"
	"go/format"
)

func main() {
	// src is the input that we want to tokenize.
	src := []byte("/** if len(sin) > 0 { */ cos(x) + 1i*$sin(x)/**}*/")
	buf := &bytes.Buffer{}
	p := newTemplE("sample.go", src, buf)
	p.run()
	fmt.Println()
	b, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("!!!", err)
		fmt.Println(string(buf.Bytes()))
	} else {
		fmt.Println(string(b))
	}

}
