package main

import (
	"mind-svc/service"
)

func main() {
	svc := service.New(service.LoadConfig())
	svc.Start()
}
