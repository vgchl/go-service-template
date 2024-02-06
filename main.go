package main

import (
	"mind-service/app"
)

func main() {
	a := app.New(app.LoadConfig())
	defer a.OnStop()
	a.Start()
}
