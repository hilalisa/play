package main

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/stores"
	"github.com/dave/play/views"
	"github.com/gopherjs/vecty"
	"github.com/vincent-petithory/dataurl"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document().(dom.HTMLDocument)

func main() {
	if document.ReadyState() == "loading" {
		document.AddEventListener("DOMContentLoaded", false, func(dom.Event) {
			go run()
		})
	} else {
		go run()
	}
}

func run() {

	vecty.AddStylesheet("https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css")
	vecty.AddStylesheet(dataurl.New([]byte(views.Styles), "text/css").String())

	app := &stores.App{}
	app.Init()
	p := views.NewPage(app)
	vecty.RenderBody(p)

	app.Dispatch(&actions.Load{})
}
