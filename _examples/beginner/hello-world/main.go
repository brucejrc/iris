package main

import (
	"github.com/brucejrc/iris"
	"github.com/brucejrc/iris/context"
)

func main() {
	app := iris.New()
	app.Handle("GET", "/", func(ctx context.Context) {
		ctx.HTML("<b> Hello world! </b>")
	})
	app.Run(iris.Addr(":8080"))
}
