package temple

import (
	"bytes"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Temple struct {
	prnt    *printer_
	scn     *scanner.Scanner
	fset    *token.FileSet
	line    int
	pos     token.Position
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
			//p.prnt.flushLine(p.pos.Line - p.line)
		}
		return token.EOF, "", true
	}

	if p.pos.Line > p.line {
		if p.started {
			//p.prnt.flushLine(p.pos.Line - p.line)
		}
		p.line = p.pos.Line
	}

	return tok, lit, false
}

func (p *Temple) errorf(format string, args ...interface{}) {
	head := fmt.Sprintf("%s:%d:%d:", p.pos.Filename, p.pos.Line, p.pos.Column)
	fmt.Printf(head+format+"\n", args...)
}

func (p *Temple) getParen() string {
	n := 1
	s := ""
	addt := func(tok token.Token, lit string) string {
		if len(lit) > 0 {
			return " " + lit
		} else {
			return " " + tok.String()
		}
	}
	for n > 0 {
		tok, lit, stop := p.scan()
		if stop {
			break
		}
		switch tok {
		case token.RPAREN:
			n -= 1
			if n > 0 {
				s += addt(tok, lit)
			}
		case token.LPAREN:
			n += 1
			s += addt(tok, lit)
		default:
			s += addt(tok, lit)
		}
	}
	return s
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
				s := lit[3:len(lit)]
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
				s := lit[3:len(lit)]
				if p.started {
					p.prnt.flush()
				}
				p.prnt.code(s)
			default:
				if !p.started {
					continue loop
				}
				p.addToken(tok, lit)
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
				case token.LPAREN:
					s := p.getParen()
					p.prnt.flush()
					p.prnt.printVar(s)
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
				case token.LPAREN:
					s := p.getParen()
					p.prnt.flush()
					p.prnt.printVarString(s)
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

func Run(args ...string) {
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
	New("sample.go", src, buf).Run()
	b, err := format.Source(buf.Bytes())
	if err != nil {
		ioutil.WriteFile(outputName, buf.Bytes(), os.ModePerm)
		log.Fatalln(err)
	}
	ioutil.WriteFile(outputName, b, os.ModePerm)
}
