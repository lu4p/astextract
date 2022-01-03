//go:build !static && !debug

package main

import (
	"fmt"

	"github.com/lu4p/astextract"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

//go:generate env GOARCH=wasm GOOS=js go build -ldflags "-s -w" -o ./web/app.wasm
//go:generate go run .

const explain = `This tool converts Go code into its go/ast representation,  using WebAssembly.`

const usage = `Paste go code (either a single ast.Expr or a whole file) on the left and the equivalent go/ast representation will be generated on the right.`

type astweb struct {
	app.Compo
	ast string
}

func (aw *astweb) Render() app.UI {
	return app.Div().Class("container").Body(
		app.H1().Text("astextract"),
		app.P().Text(explain),
		app.P().Text(usage),

		app.P().Body(
			app.A().Href("https://github.com/lu4p/astextract").Text("astextract on Github"),
			app.Text(" | "),
			app.A().Href("https://pkg.go.dev/go/ast").Text("go/ast documentation"),

			app.Text(" | Created by "),
			app.A().Href("https://github.com/lu4p").Text("lu4p"),

			app.Text(" | Created with "),
			app.A().Href("https://github.com/maxence-charriere/go-app").Text("go-app"),
		),

		app.Div().Class("columns").Body(
			app.Div().Class("column").Body(
				app.Div().Class("form-group").Body(
					app.Textarea().Rows(30).Class("form-input").Placeholder("Go code").OnInput(aw.OnChange),
				),
			),

			app.Div().Class("column").Body(
				app.Div().Class("form-group").Body(
					app.Textarea().Rows(30).Placeholder("Ast Output").Class("form-input").ReadOnly(true).Body(
						app.Text(aw.ast),
					),
				),
			),
		),
	)
}

func (aw *astweb) OnChange(ctx app.Context, e app.Event) {
	input := ctx.JSSrc().Get("value").String()

	ast, err := astextract.Parse(input)
	if err != nil {
		aw.ast = err.Error()
	} else {
		aw.ast = ast
	}

	aw.Update()
}

func main() {
	app.Route("/", &astweb{})
	app.RunWhenOnBrowser()

	h := &app.Handler{
		Name:       "astextract",
		Title:      "astextract",
		ThemeColor: "#ffffff",
		Author:     "lu4p",
		Styles:     []string{"https://unpkg.com/spectre.css/dist/spectre.min.css"},
		Resources:  app.GitHubPages("astextract"),
	}

	err := app.GenerateStaticWebsite("", h)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("static website generated")
}
