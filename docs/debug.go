// +build debug

package main

import (
	"fmt"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	h := &app.Handler{
		Name:       "astextract",
		Title:      "astextract",
		Author:     "lu4p",
		ThemeColor: "#ffffff",
		Styles:     []string{"https://unpkg.com/spectre.css/dist/spectre.min.css"},
	}

	fmt.Println("listening on http://localhost:7778")
	if err := http.ListenAndServe(":7778", h); err != nil {
		panic(err)
	}
}
