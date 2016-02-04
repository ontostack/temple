package main

import (
	"flag"
	"fmt"
	"go/scanner"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/benbjohnson/ego"
)

func run(out, pkg, ext string, roots []string) {
	// Recursively retrieve all templates
	v := visitor{ext: ext}
	for _, root := range roots {
		if err := filepath.Walk(root, v.visit); err != nil {
			scanner.PrintError(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Parse every template file.
	var templates []*ego.Template
	for _, path := range v.paths {
		t, err := ego.ParseFile(path)
		if err != nil {
			log.Fatal("parse file: ", err)
		}
		templates = append(templates, t)
	}

	// If we have no templates then exit.
	if len(templates) == 0 {
		os.Exit(0)
	}

	// Write package to output file.
	p := &ego.Package{Templates: templates, Name: pkg}
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Write template to file.
	if err := p.Write(f); err != nil {
		log.Fatal("write: ", err)
	}
}

func main() {
	outfile := flag.String("o", "templ.go", "output file")
	pkgname := flag.String("p", "", "package name")
	ext := flag.String("e", "go", "template file extension")
	flag.Parse()
	log.SetFlags(0)

	// If no paths are provided then use the present working directory.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// If no package name is set then use the directory name of the output file.
	if *pkgname == "" {
		abspath, _ := filepath.Abs(*outfile)
		*pkgname = filepath.Base(filepath.Dir(abspath))
		*pkgname = regexp.MustCompile(`(\w+).*`).ReplaceAllString(*pkgname, "$1")
	}

	run(*outfile, *pkgname, *ext, roots)
}

// visitor iterates over
type visitor struct {
	paths []string
	ext string
}

func (v *visitor) visit(path string, info os.FileInfo, err error) error {
	if info == nil {
		return fmt.Errorf("file not found: %s", path)
	}
	if !info.IsDir() && strings.HasSuffix(filepath.Base(path), "." + v.ext) {
		v.paths = append(v.paths, path)
	}
	return nil
}
