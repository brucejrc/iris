package main

import (
	"github.com/brucejrc/iris"
	"github.com/brucejrc/iris/context"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx context.Context) {
		file := "./files/first.zip"
		ctx.SendFile(file, "c.zip")
	})

	app.Run(iris.Addr(":8080"))
}
