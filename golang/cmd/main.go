package main

import (
	"testtask/internal/app"
	"testtask/internal/config"
)

func main() {
	cfg := config.MustNew()
	app.Run(cfg)
}
