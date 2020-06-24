package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func genFile(outfile string, b *bytes.Buffer) error {
	y := `package main

	import (
		"go/ast"
		"go/printer"
		"go/token"
		"log"
		"os"
	)
	
	func main() {
		z := ` + b.String() + `

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

	return ioutil.WriteFile(outfile, []byte(y), 0644)
}
