package main

import (
	"mind-svc/app"
)

func main() {
	a := app.New(app.LoadConfig())
	a.Start()
}
