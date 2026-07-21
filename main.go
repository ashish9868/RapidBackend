package main

import (
	"embed"

	"github.com/ashish9868/rapidbackend/cmd"
	"github.com/ashish9868/rapidbackend/core"
)

var (
	VERSION = "v0.0.1-prod"
)

//go:embed frontend/dist/index.html
//go:embed frontend/dist/**/*
var embeddedFiles embed.FS

func main() {
	app := core.NewApp(embeddedFiles)
	app.Version = &VERSION
	cmd.ExecuteRootCommand(app)
}
