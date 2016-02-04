package main

import (
	"go/format"
	"bytes"
	"strings"
	"fmt"
	"go/scanner"
	"go/token"
)

type printer struct {
	scn  *scanner.Scanner
	data string
	line int
	fset *token.FileSet
	buf *bytes.Buffer
}

func newPrinter(src []byte) *printer {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)
	return &printer{
		scn: &s,
		fset: fset,
		buf: &bytes.Buffer{},
	}
}

func (p *printer) flush() {
	if len(p.data) > 0 {
		fmt.Fprintf(p.buf, "fmt.Print(`%s`)\n", p.data)
		p.data = ""
	}
}

func (p *printer) flushLine(inc int) {
	p.flush()
	s := ""
	for i := 0; i < inc; i++ {
		s += `\n`
	}
	if inc > 0 {
		fmt.Fprintf(p.buf, `fmt.Print("%s")` + "\n", s)
	}
}

func (p *printer) addToken(tok token.Token, lit string) {
	if len(lit) > 0 {
		p.data += " " + lit
	} else {
		p.data += " " + tok.String()
	}

	fmt.Printf(">><%s, %s>: `%s`\n", tok.String(), lit, p.data)
}

func (p *printer) scan()(token.Token, string, bool) {
	fpos, tok, lit := p.scn.Scan()
	pos := p.fset.Position(fpos)
	if tok == token.EOF {
		p.flushLine(pos.Line-p.line)
		return token.EOF, "", true
	}

	if pos.Line > p.line {
		p.flushLine(pos.Line-p.line)
		p.line = pos.Line
	}

	return tok, lit, false
}

func (p *printer) run() {
	loop:
	for {
		tok, lit, stop := p.scan()
		if stop {
			break loop
		}

		switch {
		case tok == token.COMMENT:
			if strings.HasPrefix(lit, "/**") {
				s := lit[3:len(lit)-2]
				p.flush()
				fmt.Fprintln(p.buf, s)
			}
		case tok == token.ILLEGAL:
			if lit == "$" {
				tok, lit, stop := p.scan()
				if tok != token.IDENT {
					fmt.Println("!!!Unexpected token: " + tok.String())
					return
				}
				if stop {
					break loop
				}
				p.flush()
				fmt.Fprintf(p.buf, "fmt.Print(%s)\n", lit)
			}
		default:
			p.addToken(tok, lit)
		}
	}
}

func main() {
	// src is the input that we want to tokenize.
	src := []byte("/** if len(sin) > 0 { */ cos(x) + 1i*$sin(x)/**}*/")
	p := newPrinter(src)
	p.run()
	fmt.Println()
	b, err := format.Source(p.buf.Bytes())
	if err != nil {
		fmt.Println("!!!", err)
	} else {
		fmt.Println(string(b))
	}
}
