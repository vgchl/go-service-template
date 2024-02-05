package main

import (
	"mind-service/app"
)

func main() {
	a := app.New(app.LoadConfig())
	a.Start()
}
