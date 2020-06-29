// +build !debug

package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {

	h := &app.Handler{
		Name:      "astextract",
		Title:     "astextract",
		Author:    "lu4p",
		Styles:    []string{"https://unpkg.com/spectre.css/dist/spectre.min.css"},
		Resources: app.GitHubPages("astextract"),
	}

	err := app.GenerateStaticWebsite("../static", h)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("static website generated")

}
