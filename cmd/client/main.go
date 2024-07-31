package main

import (
	"github.com/nextlag/keeper/internal/client/app"
	"github.com/nextlag/keeper/internal/client/app/build"
)

func main() {
	build.Version = "1.0.0"
	build.CheckBuild()
	app.Execute()
}
