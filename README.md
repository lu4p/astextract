# astextract
[![Test](https://github.com/lu4p/astextract/workflows/Test/badge.svg)](https://github.com/lu4p/astextract/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/lu4p/astextract)](https://goreportcard.com/report/github.com/lu4p/astextract)


astextract converts a given go file to its [ast](https://pkg.go.dev/go/ast) representation.

This is useful for easliy writing typesafe `go generate` tools, which don't concatenate strings for generating code. 

The output of astextract is 100% valid go, so it can be used in go code directly without any modifications.

All zero/null value struct fields are ommited for a more compact ast representation.

Absolute Position info is stripped, because positions change if you add any dynamic content to the ast.

I already wrote [binclude](https://github.com/lu4p/binclude) a tool for resource embedding with the help of astextraxt, and added string obfuscation to [garble](https://github.com/mvdan/garble).

## Web app
astextract be used as a Web app: https://astextract.lu4p.xyz/

The Web app was created using the amazing [go-app](https://github.com/maxence-charriere/go-app) package.


## Install
`GO111MODULE=on go get -u github.com/lu4p/astextract/cmd/astextract`

## Usage
`astextract [flags] file`

See `astextract -h` for up to date usage information.

## Example 
main.go:
```go
package main

func main() {
	println("Hello, World!")
}
```

Command: `astextract main.go`

Output:
```go
&ast.File {
  Package: 1,
  Name: &ast.Ident {
    Name: "main",
  },
  Decls: []ast.Decl {
    &ast.FuncDecl {
      Name: &ast.Ident {
        Name: "main",
      },
      Type: &ast.FuncType {
        Params: &ast.FieldList {},
      },
      Body: &ast.BlockStmt {
        List: []ast.Stmt {
          &ast.ExprStmt {
            X: &ast.CallExpr {
              Fun: &ast.Ident {
                Name: "println",
              },
              Args: []ast.Expr {
                &ast.BasicLit {
                  Kind: token.STRING,
                  Value: "\"Hello, World!\"",
                },
              },
            },
          },
        },
      },
    },
  },
}
```


## How to convert go/ast back to go
Example taken from binclude [generateFile()](https://github.com/lu4p/binclude/blob/f796cc285c1e76e072386a639b8209284aaf3369/cmd/binclude/inject.go#L160)

```go
bincludeFile :=  &ast.File{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Slash: 1,
					Text:  "// Code generated by binclude; DO NOT EDIT.",
				},
			},
		},
		Package: 45,
		Name:    pkgName,
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: imports,
			},
			&ast.GenDecl{
				Tok:   token.VAR,
				Specs: astVars,
			},
		},
    }

f, err := os.OpenFile("binclude.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
if err != nil {
	return err
}
defer f.Close()

err = printer.Fprint(f, fset, bincludeFile)
if err != nil {
	return err
}
```

`astextract -out=filename main.go` generates a go file containing the ast of `main.go` and an example on how to generate go code from the ast.

## Alternatives
- [go2ast](https://github.com/reflog/go2ast): only converts a single line and doesn't support top level definitions.
- [go/ast Print](https://pkg.go.dev/go/ast?tab=doc#Print): prints a representation which is unsuitable for usage in code generation

From the go/ast documentation:

This example shows what an AST looks like when printed for debugging. 
Code:
```go
// src is the input for which we want to print the AST.
src := `
package main
func main() {
	println("Hello, World!")
}
`

// Create the AST by parsing src.
fset := token.NewFileSet() // positions are relative to fset
f, err := parser.ParseFile(fset, "", src, 0)
if err != nil {
	panic(err)
}

// Print the AST.
ast.Print(fset, f)
```
Output: 
```
     0  *ast.File {
     1  .  Package: 2:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 2:9
     4  .  .  Name: "main"
     5  .  }
     6  .  Decls: []ast.Decl (len = 1) {
     7  .  .  0: *ast.FuncDecl {
     8  .  .  .  Name: *ast.Ident {
     9  .  .  .  .  NamePos: 3:6
    10  .  .  .  .  Name: "main"
    11  .  .  .  .  Obj: *ast.Object {
    12  .  .  .  .  .  Kind: func
    13  .  .  .  .  .  Name: "main"
    14  .  .  .  .  .  Decl: *(obj @ 7)
    15  .  .  .  .  }
    16  .  .  .  }
    17  .  .  .  Type: *ast.FuncType {
    18  .  .  .  .  Func: 3:1
    19  .  .  .  .  Params: *ast.FieldList {
    20  .  .  .  .  .  Opening: 3:10
    21  .  .  .  .  .  Closing: 3:11
    22  .  .  .  .  }
    23  .  .  .  }
    24  .  .  .  Body: *ast.BlockStmt {
    25  .  .  .  .  Lbrace: 3:13
    26  .  .  .  .  List: []ast.Stmt (len = 1) {
    27  .  .  .  .  .  0: *ast.ExprStmt {
    28  .  .  .  .  .  .  X: *ast.CallExpr {
    29  .  .  .  .  .  .  .  Fun: *ast.Ident {
    30  .  .  .  .  .  .  .  .  NamePos: 4:2
    31  .  .  .  .  .  .  .  .  Name: "println"
    32  .  .  .  .  .  .  .  }
    33  .  .  .  .  .  .  .  Lparen: 4:9
    34  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    35  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    36  .  .  .  .  .  .  .  .  .  ValuePos: 4:10
    37  .  .  .  .  .  .  .  .  .  Kind: STRING
    38  .  .  .  .  .  .  .  .  .  Value: "\"Hello, World!\""
    39  .  .  .  .  .  .  .  .  }
    40  .  .  .  .  .  .  .  }
    41  .  .  .  .  .  .  .  Ellipsis: -
    42  .  .  .  .  .  .  .  Rparen: 4:25
    43  .  .  .  .  .  .  }
    44  .  .  .  .  .  }
    45  .  .  .  .  }
    46  .  .  .  .  Rbrace: 5:1
    47  .  .  .  }
    48  .  .  }
    49  .  }
    50  .  Scope: *ast.Scope {
    51  .  .  Objects: map[string]*ast.Object (len = 1) {
    52  .  .  .  "main": *(obj @ 11)
    53  .  .  }
    54  .  }
    55  .  Unresolved: []*ast.Ident (len = 1) {
    56  .  .  0: *(obj @ 29)
    57  .  }
    58  }
```

