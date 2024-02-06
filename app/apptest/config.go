package apptest

import "mind-service/app"

func Config() app.Config {
	c := app.DefaultConfig()
	c.Env = "local"
	c.LogLevel = "debug"
	return c
}
