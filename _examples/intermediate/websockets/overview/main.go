package main

import (
	"fmt" // optional

	"github.com/brucejrc/iris"
	"github.com/brucejrc/iris/context"

	"github.com/brucejrc/iris/view"
	"github.com/brucejrc/iris/websocket"
)

type clientPage struct {
	Title string
	Host  string
}

func main() {
	app := iris.New()
	app.AttachView(view.HTML("./templates", ".html")) // select the html engine to serve templates

	ws := websocket.New(websocket.Config{
		// the request path which the websocket client should listen/registered to the server,
		// the endpoint is like a route.
		Endpoint: "/my_endpoint",
		// the client-side javascript static file path
		// which will be served by Iris.
		// default is /iris-ws.js
		// if you change that you have to change the bottom of templates/client.html
		// script tag:
		ClientSourcePath: "/iris-ws.js",
		//
		// Set the timeouts, 0 means no timeout
		// websocket has more configuration, go to ../../config.go for more:
		// WriteTimeout: 0,
		// ReadTimeout:  0,
		// by-default all origins are accepted, you can change this behavior by setting:
		// CheckOrigin: (r *http.Request ) bool {},
		//
		//
		// IDGenerator used to create (and later on, set)
		// an ID for each incoming websocket connections (clients).
		// The request is an argument which you can use to generate the ID (from headers for example).
		// If empty then the ID is generated by DefaultIDGenerator: randomString(64):
		// IDGenerator func(ctx context.Context) string {},
	})

	ws.Attach(app) // adapt the websocket server, you can adapt more than one with different Endpoint

	app.StaticWeb("/js", "./static/js") // serve our custom javascript code

	app.Get("/", func(ctx context.Context) {
		ctx.ViewData("", clientPage{"Client Page", ctx.Host()})
		ctx.View("client.html")
	})

	var myChatRoom = "room1"

	ws.OnConnection(func(c websocket.Connection) {
		// Context returns the (upgraded) context.Context of this connection
		// avoid using it, you normally don't need it,
		// websocket has everything you need to authenticate the user BUT if it's necessary
		// then  you use it to receive user information, for example: from headers.

		// ctx := c.Context()

		// join to a room (optional)
		c.Join(myChatRoom)

		c.On("chat", func(message string) {
			if message == "leave" {
				c.Leave(myChatRoom)
				c.To(myChatRoom).Emit("chat", "Client with ID: "+c.ID()+" left from the room and cannot send or receive message to/from this room.")
				c.Emit("chat", "You have left from the room: "+myChatRoom+" you cannot send or receive any messages from others inside that room.")
				return
			}
			// to all except this connection ->
			// c.To(websocket.Broadcast).Emit("chat", "Message from: "+c.ID()+"-> "+message)
			// to all connected clients: c.To(websocket.All)

			// to the client itself ->
			//c.Emit("chat", "Message from myself: "+message)

			//send the message to the whole room,
			//all connections are inside this room will receive this message
			c.To(myChatRoom).Emit("chat", "From: "+c.ID()+": "+message)
		})

		// or create a new leave event
		// c.On("leave", func() {
		// 	c.Leave(myChatRoom)
		// })

		c.OnDisconnect(func() {
			fmt.Printf("Connection with ID: %s has been disconnected!\n", c.ID())
		})
	})

	app.Run(iris.Addr(":8080"))
}
