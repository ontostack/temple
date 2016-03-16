package temple

import (
	"fmt"
	"io"
)

type printer_ struct {
	wr    io.Writer
	data  string
	wname string
}

func (p *printer_) print(s string) {
	fmt.Fprintf(p.wr, p.wname+".Write([]byte(`%s`))\n", s)
}

func (p *printer_) println(s string) {
	fmt.Fprintf(p.wr, p.wname+".Write([]byte(`%s`+"+`"\n"))`+"\n", s)
}

func (p *printer_) code(s string) {
	fmt.Fprint(p.wr, s)
}

func (p *printer_) printVar(s string) {
	fmt.Fprintf(p.wr, p.wname+".Write([]byte(%s))\n", s)
}

func (p *printer_) printVarString(s string) {
	fmt.Fprintf(p.wr, p.wname+".Write([]byte(`\"`+%s+`\"`))\n", s)
}

func (p *printer_) addData(s string) {
	p.data += s
}

func (p *printer_) flush() {
	if len(p.data) > 0 {
		p.print(p.data)
		p.data = ""
	}
}

func (p *printer_) flushLine(lines int) {
	p.flush()
	for i := 0; i < lines; i++ {
		p.println("")
	}
}
