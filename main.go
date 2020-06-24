package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
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

func main() {
	os.Exit(main1())
}

func main1() int {
	flag.Parse()

	path := flag.Arg(0)

	if path == "" {
		flag.Usage()
	}

	err := parseFile(path)
	if err != nil {
		log.Println(err)
		return 1
	}

	return 0
}

func parseFile(path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	filters := AppendFilters(PosFilter, KeywordFilter, ZeroFilter)

	if outfile != "" {

		var b bytes.Buffer

		err = Fprint(&b, fset, f, filters)
		if err != nil {
			return err
		}

		return genFile(outfile, &b)
	}

	return Fprint(os.Stdout, fset, f, filters)
}
