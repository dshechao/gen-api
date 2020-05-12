// file: main.go
package main

import (
	"github.com/kataras/iris/v12"

	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/irisgen"
)

type myXML struct {
	Result string `xml:"result"`
}

type myModel struct {
	Username string `json:"username"`
	Gender   string `json:"gender"`
}

// create a function to initialize our app and gen middleware
// in order to be used on the test file as well.
func newApp() *iris.Application {
	app := iris.New()

	gen.Init(&gen.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: "Iris",
		DocPath:  "doc/index.html",
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})

	app.Use(irisgen.New()) // <- IMPORTANT, register the middleware.

	app.Get("/json", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"result": "Hello World!"})
	})

	app.Get("/plain", func(ctx iris.Context) {
		ctx.Text("Hello World!")
	})

	app.Get("/xml", func(ctx iris.Context) {
		ctx.XML(myXML{Result: "Hello World!"})
	})

	app.Get("/complex", func(ctx iris.Context) {
		value := ctx.URLParam("key")
		ctx.JSON(iris.Map{"value": value})
	})

	app.Post("/reqbody", func(ctx iris.Context) {
		var model myModel
		ctx.ReadJSON(&model)
		ctx.Writef(model.Username)
	})

	app.Post("/hello", func(ctx iris.Context) {
		username := ctx.FormValue("username")
		ctx.Writef("Hello %s", username)
	})

	return app
}

func main() {
	app := newApp()
	// Run our HTTP Server.
	//
	// Note that on each incoming request the gen will generate and update the "apidoc.html".
	// Recommentation:
	// Write tests that calls those handlers, save the generated "apidoc.html".
	// Turn off the gen middleware when in production.
	//
	// Example usage:
	// Visit all paths and open the generated "apidoc.html" file to see the API's automated docs.
	app.Run(iris.Addr(":8000"))
}
