package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"strings"
)

type temple struct {
	prnt *printer_
	scn  *scanner.Scanner
	fset *token.FileSet
	line int
	pos  token.Position
}

func newTemplE(fname string, src []byte, wr io.Writer) *temple {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(fname, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)
	return &temple{
		scn:  &s,
		fset: fset,
		prnt: &printer_{wr: wr},
	}
}

func (p *temple) addToken(tok token.Token, lit string) {
	if len(lit) > 0 {
		p.prnt.addData(" " + lit)
	} else {
		p.prnt.addData(" " + tok.String())
	}
}

func (p *temple) scan() (token.Token, string, bool) {
	fpos, tok, lit := p.scn.Scan()
	p.pos = p.fset.Position(fpos)
	if tok == token.EOF {
		p.prnt.flushLine(p.pos.Line - p.line)
		return token.EOF, "", true
	}

	if p.pos.Line > p.line {
		p.prnt.flushLine(p.pos.Line - p.line)
		p.line = p.pos.Line
	}

	return tok, lit, false
}

func (p *temple) errorf(format string, args ...interface{}) {
	head := fmt.Sprintf("%s:%d:%d:", p.pos.Filename, p.pos.Line, p.pos.Column)
	fmt.Printf(head+format+"\n", args...)
}

func (p *temple) run() {
loop:
	for {
		tok, lit, stop := p.scan()
		if stop {
			break loop
		}

		switch {
		case tok == token.COMMENT:
			switch {
			case strings.HasPrefix(lit, "/**"):
				s := lit[3 : len(lit)-2]
				p.prnt.flush()
				p.prnt.code(s)
			case strings.HasPrefix(lit, "///"):
				s := lit[3 : len(lit)-1]
				p.prnt.flush()
				p.prnt.code(s)
			}
		case tok == token.ILLEGAL:
			if lit == "$" {
				tok, lit, stop := p.scan()
				if stop {
					break loop
				}
				switch tok {
				case token.IDENT:
					p.prnt.flush()
					p.prnt.printVar(lit)
				default:
					p.errorf("Unexpected token: %s", tok.String())
					return
				}
			}
		default:
			p.addToken(tok, lit)
		}
	}
}
