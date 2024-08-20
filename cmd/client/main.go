package main

import (
	"github.com/nextlag/keeper/internal/client/app"
	"github.com/nextlag/keeper/internal/client/app/build"
)

func main() {
	build.CheckBuild()
	app.Execute()
}
