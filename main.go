package main

import (
	c "github.com/yuxuan0105/gin_practice/controller"
)

func main() {
	app := c.NewController()

	app.RunServer()
}
