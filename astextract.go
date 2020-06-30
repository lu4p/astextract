package astextract

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

var outfile string

func init() {
	flag.Usage = usage
	flag.StringVar(&outfile, "out", "", "write output to a file instead of stdout")
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage of astextract:

astextract [flags] file

astextract accepts the following flags:
`)
	flag.PrintDefaults()
	os.Exit(2)
}

// Main1 called by cmd/astextract
func Main1() int {
	flag.Parse()

	path := flag.Arg(0)

	if path == "" {
		flag.Usage()
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return 1
	}

	out, err := Parse(string(content))
	if err != nil {
		log.Println(err)
		return 1
	}

	if outfile != "" {
		err := genFile(outfile, out)
		if err != nil {
			log.Println(err)
			return 1
		}

		return 0
	}

	fmt.Println(out)

	return 0
}

// Parse converts input to ast
func Parse(input string) (string, error) {
	var b bytes.Buffer

	fset := token.NewFileSet()
	filters := AppendFilters(PosFilter, KeywordFilter, ZeroFilter)

	expr, err := parser.ParseExpr(input)
	if err == nil {
		err = Fprint(&b, fset, expr, filters)
		return b.String(), err
	}

	f, err := parser.ParseFile(fset, "", input, parser.ParseComments)
	if err != nil {
		return "", err
	}

	for i, comment := range f.Comments {
		if comment.Pos() > f.Package {
			f.Comments = f.Comments[:i]
		}
	}

	err = Fprint(&b, fset, f, filters)
	if err != nil {
		return "", err
	}

	return b.String(), err
}
