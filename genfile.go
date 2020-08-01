package astextract

import (
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
)

func genFile(outfile string, out string) error {
	y := `package main

import (
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
)
	
func main() {
	z := ` + out + `

	// example usage
	f, err := os.OpenFile("gen.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	fSet := token.NewFileSet()

	printer.Fprint(f, fSet, z)
}`

	err := os.MkdirAll(filepath.Dir(outfile), 0755)
	if err != nil {
		return err
	}

	formatted, err := format.Source([]byte(y))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outfile, formatted, 0644)
}
