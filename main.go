package main

import (
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
)

func main() {
	app, err := v1.NewModel()
	if err != nil {
		panic(err)
	}

	app.RunServer()
}
