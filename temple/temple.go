package temple

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"strings"
)

type Temple struct {
	prnt *printer_
	scn  *scanner.Scanner
	fset *token.FileSet
	line int
	pos  token.Position
	started bool
}

func New(fname string, src []byte, wr io.Writer) *Temple {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(fname, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)
	return &Temple{
		scn:  &s,
		fset: fset,
		prnt: &printer_{wr: wr},
	}
}

func (p *Temple) addToken(tok token.Token, lit string) {
	if len(lit) > 0 {
		p.prnt.addData(" " + lit)
	} else {
		p.prnt.addData(" " + tok.String())
	}
}

func (p *Temple) scan() (token.Token, string, bool) {
	fpos, tok, lit := p.scn.Scan()
	p.pos = p.fset.Position(fpos)
	if tok == token.EOF {
		if p.started {
			p.prnt.flushLine(p.pos.Line - p.line)
		}
		return token.EOF, "", true
	}

	if p.pos.Line > p.line {
		if p.started {
			p.prnt.flushLine(p.pos.Line - p.line)
		}
		p.line = p.pos.Line
	}

	return tok, lit, false
}

func (p *Temple) errorf(format string, args ...interface{}) {
	head := fmt.Sprintf("%s:%d:%d:", p.pos.Filename, p.pos.Line, p.pos.Column)
	fmt.Printf(head+format+"\n", args...)
}

func (p *Temple) Run() {
	p.started = false
loop:
	for {
		tok, lit, stop := p.scan()
		if stop {
			break loop
		}

		switch {
		case tok == token.COMMENT:
			switch {
			case strings.HasPrefix(lit, "/*-"):
				s := lit[3 : len(lit)-2]
				if p.started {
					p.prnt.flush()
				}
				p.prnt.code(s)
				p.started = !p.started
			case strings.HasPrefix(lit, "//-"):
				s := lit[3 : len(lit)]
				if p.started {
					p.prnt.flush()
				}
				p.prnt.code(s)
				p.started = !p.started
			case strings.HasPrefix(lit, "/**"):
				s := lit[3 : len(lit)-2]
				if p.started {
					p.prnt.flush()
				}
				p.prnt.code(s)
			case strings.HasPrefix(lit, "///"):
				s := lit[3 : len(lit)]
				if p.started {
					p.prnt.flush()
				}
				p.prnt.code(s)
			}
		case tok == token.ILLEGAL:
			if !p.started {
				continue loop
			}
			switch lit {
			case "$":
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
			case "#":
				tok, lit, stop := p.scan()
				if stop {
					break loop
				}
				switch tok {
				case token.IDENT:
					p.prnt.flush()
					p.prnt.printVarString(lit)
				default:
					p.errorf("Unexpected token: %s", tok.String())
					return
				}
			}
		default:
			if !p.started {
				continue loop
			}
			p.addToken(tok, lit)
		}
	}
}
