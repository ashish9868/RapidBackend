package main

import (
	"embed"

	"github.com/ashish9868/rapidbackend/cmd"
	"github.com/ashish9868/rapidbackend/core"
)

//go:embed templates/*
//go:embed templates/**/*
var embeddedFiles embed.FS

func main() {
	app := core.NewApp(embeddedFiles)
	cmd.ExecuteRootCommand(app)
}
