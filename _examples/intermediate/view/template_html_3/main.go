// Package main an example on how to naming your routes & use the custom 'url path' HTML Template Engine, same for other template engines.
package main

import (
	"github.com/brucejrc/iris"
	"github.com/brucejrc/iris/context"
	"github.com/brucejrc/iris/view"
)

func main() {
	app := iris.New()
	if err := app.AttachView(view.HTML("./templates", ".html").Reload(true)); err != nil {
		panic(err)
	}

	mypathRoute, _ := app.Get("/mypath", writePathHandler)
	mypathRoute.Name = "my-page1"

	mypath2Route, err := app.Get("/mypath2/{paramfirst}/{paramsecond}", writePathHandler)
	// same as: app.Get("/mypath2/:paramfirst/:paramsecond", writePathHandler)
	if err != nil { // catch errors when creating a route or catch them on err := .Run, it's up to you.
		panic(err)
	}
	mypath2Route.Name = "my-page2"

	mypath3Route, _ := app.Get("/mypath3/{paramfirst}/statichere/{paramsecond}", writePathHandler)
	mypath3Route.Name = "my-page3"

	mypath4Route, _ := app.Get("/mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}", writePathHandler)
	// same as: app.Get("/mypath4/:paramfirst/statichere/:paramsecond/:otherparam/*something", writePathHandler)
	mypath4Route.Name = "my-page4"

	// same with Handle/Func
	mypath5Route, _ := app.Handle("GET", "/mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{something:path}", writePathHandler)
	mypath5Route.Name = "my-page5"

	mypath6Route, _ := app.Get("/mypath6/{paramfirst}/{paramsecond}/statichere/{paramThirdAfterStatic}", writePathHandler)
	mypath6Route.Name = "my-page6"

	app.Get("/", func(ctx context.Context) {
		// for /mypath6...
		paramsAsArray := []string{"theParam1", "theParam2", "paramThirdAfterStatic"}
		ctx.ViewData("ParamsAsArray", paramsAsArray)
		if err := ctx.View("page.html"); err != nil {
			panic(err)
		}
	})

	app.Get("/redirect/{namedRoute}", func(ctx context.Context) {
		routeName := ctx.Params().Get("namedRoute")
		r := app.GetRoute(routeName)
		if r == nil {
			ctx.StatusCode(404)
			ctx.Writef("Route with name %s not found", routeName)
			return
		}

		println("The path of " + routeName + "is: " + r.Path)
		// if routeName == "my-page1"
		// prints: The path of of my-page1 is: /mypath
		// if it's a path which takes named parameters
		// then use "r.ResolvePath(paramValuesHere)"
		ctx.Redirect(r.Path)
		// http://localhost:8080/redirect/my-page1 will redirect to -> http://localhost:8080/mypath
	})

	// http://localhost:8080
	// http://localhost/redirect/my-page1
	app.Run(iris.Addr(":8080"))

}

func writePathHandler(ctx context.Context) {
	ctx.Writef("Hello from %s.", ctx.Path())
}
