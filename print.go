// This file is adapted from https://golang.org/src/go/ast/print.go to print ast in go syntax compatible form

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains printing support for ASTs.

package astextract

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strings"
)

var skipBlock = false

// A FieldFilter may be provided to Fprint to control the output.
type FieldFilter func(name string, value reflect.Value) bool

// ZeroFilter returns false for field values that are the zero value
func ZeroFilter(name string, v reflect.Value) bool {
	if skipBlock {
		return true
	}
	return !v.IsZero()
}

var types = []string{
	"ArrayType", "AssignStmt", "BasicLit", "BinaryExpr", "BlockStmt",
	"BranchStmt", "CallExpr", "CaseClause", "ChanType",
	"CommClause", "CompositeLit", "DeclStmt", "DeferStmt",
	"Ellipsis", "ExprStmt", "Field", "FieldList", "File",
	"ForStmt", "FuncDecl", "FuncLit", "FuncType", "GenDecl", "GoStmt",
	"Ident", "IfStmt", "ImportSpec", "IncDecStmt", "IndexExpr", "InterfaceType",
	"KeyValueExpr", "LabeledStmt", "MapType", "ParenExpr", "RangeStmt",
	"ReturnStmt", "SelectStmt", "SelectorExpr", "SendStmt", "SliceExpr",
	"StarExpr", "StructType", "SwitchStmt", "TypeAssertExpr", "TypeSpec",
	"TypeSwitchStmt", "UnaryExpr", "ValueSpec",
}

// UnusedTypes this is for checking test coverage
var UnusedTypes map[string]uint8

func init() {
	genUnusedTypes()
}

func genUnusedTypes() {
	UnusedTypes = make(map[string]uint8, len(types))
	for _, t := range types {
		UnusedTypes[t] = 0
	}
}

// PosFilter returns false for field names which are Positions
func PosFilter(name string, v reflect.Value) bool {
	if strings.HasSuffix(name, "Pos") {
		return false
	}

	blacklist := []string{
		"If", "Return", "Func", "Opening", "Closing", "Colon", "Obj",
		"Struct", "Map", "For", "Star", "Case", "Begin", "Defer", "Go",
		"Interface", "Select", "Struct", "Switch", "Arrow",
	}

	bracePos := []string{"brace", "paren", "brack"}
	for _, item := range bracePos {
		blacklist = append(blacklist, "L"+item, "R"+item)
	}

	for _, item := range blacklist {
		if name == item {
			return false
		}
	}

	return true
}

// KeywordFilter returns false for blacklisted field names
func KeywordFilter(name string, v reflect.Value) bool {
	var blacklist = []string{
		"Scope", "Unresolved", "Doc",
	}

	for _, item := range blacklist {
		if name == item {
			return false
		}
	}

	return true
}

// AppendFilters use multiple filters at once
func AppendFilters(filters ...FieldFilter) FieldFilter {
	return func(name string, value reflect.Value) bool {
		for _, filter := range filters {
			if !filter(name, value) {
				return false
			}
		}
		return true
	}
}

// Fprint prints the (sub-)tree starting at AST node x to w.
// If fset != nil, position information is interpreted relative
// to that file set. Otherwise positions are printed as integer
// values (file set specific offsets).
//
// A non-nil FieldFilter f may be provided to control the output:
// struct fields for which f(fieldname, fieldvalue) is true are
// printed; all others are filtered from the output. Unexported
// struct fields are never printed.
func Fprint(w io.Writer, fset *token.FileSet, x any, f FieldFilter) error {
	return fprint(w, fset, x, f)
}

func fprint(w io.Writer, fset *token.FileSet, x any, f FieldFilter) (err error) {
	// setup printer
	p := printer{
		output: w,
		fset:   fset,
		filter: f,
		ptrmap: make(map[any]int),
		last:   '\n', // force printing of line number on first line
	}

	// install error handler
	defer func() {
		if e := recover(); e != nil {
			err = e.(localError).err // re-panics if it's not a localError
		}
	}()

	// print x
	if x == nil {
		p.printf("nil\n")
		return
	}
	p.print(reflect.ValueOf(x))
	p.printf("\n")

	return
}

type printer struct {
	output io.Writer
	fset   *token.FileSet
	filter FieldFilter
	ptrmap map[any]int // *T -> line number
	indent int         // current indentation level
	last   byte        // the last byte processed by Write
	line   int         // current line number
}

var indent = []byte("  ")

func (p *printer) Write(data []byte) (n int, err error) {
	var m int
	for i, b := range data {
		// invariant: data[0:n] has been written
		if b == '\n' {
			m, err = p.output.Write(data[n : i+1])
			n += m
			if err != nil {
				return
			}
			p.line++
		} else if p.last == '\n' {
			for j := p.indent; j > 0; j-- {
				_, err = p.output.Write(indent)
				if err != nil {
					return
				}
			}
		}
		p.last = b
	}
	if len(data) > n {
		m, err = p.output.Write(data[n:])
		n += m
	}
	return
}

// localError wraps locally caught errors so we can distinguish
// them from genuine panics which we don't want to return as errors.
type localError struct {
	err error
}

// printf is a convenience wrapper that takes care of print errors.
func (p *printer) printf(format string, args ...any) {
	if _, err := fmt.Fprintf(p, format, args...); err != nil {
		panic(localError{err})
	}
}

// Implementation note: Print is written for AST nodes but could be
// used to print arbitrary data structures; such a version should
// probably be in a different package.
//
// Note: This code detects (some) cycles created via pointers but
// not cycles that are created via slices or maps containing the
// same slice or map. Code for general data structures probably
// should catch those as well.

func (p *printer) print(x reflect.Value) {
	switch x.Kind() {
	case reflect.Interface:
		p.print(x.Elem())

	case reflect.Map:
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for _, key := range x.MapKeys() {
				p.print(key)
				p.printf(": ")
				p.print(x.MapIndex(key))
				p.printf(",\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Ptr:
		p.printf("&")
		// type-checked ASTs may contain cycles - use ptrmap
		// to keep track of objects that have been printed
		// already and print the respective line number instead
		ptr := x.Interface()

		p.ptrmap[ptr] = p.line
		p.print(x.Elem())

	case reflect.Array:
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf(",\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Slice:
		if s, ok := x.Interface().([]byte); ok {
			p.printf("%#q", s)
			return
		}
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.print(x.Index(i))
				p.printf(",\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Struct:
		t := x.Type()
		delete(UnusedTypes, t.Name())

		p.printf("%s {", t)
		p.indent++
		first := true
		for i, n := 0, t.NumField(); i < n; i++ {
			// exclude non-exported fields because their
			// values cannot be accessed via reflection
			if name := t.Field(i).Name; ast.IsExported(name) {
				value := x.Field(i)
				if p.filter == nil || p.filter(name, value) {
					if p.indent == 1 {
						skipBlock = name == "Name" // skip package name
					}
					if first {
						p.printf("\n")
						first = false
					}
					p.printf("%s: ", name)
					p.print(value)
					p.printf(",\n")
				}
			}
		}
		p.indent--
		p.printf("}")

	default:
		v := x.Interface()
		switch v := v.(type) {
		case string:
			// print strings in quotes
			p.printf("%q", v)
		case token.Token:
			p.printf("token.%v", tokenString(v))
		default:
			p.printf("%v", v)
		}
	}
}
