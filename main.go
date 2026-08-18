package main

import (
	"embed"

	"github.com/ashish9868/rapidbackend/cmd"
	"github.com/ashish9868/rapidbackend/core"
)

var (
	VERSION = "v0.0.1-prod"
)

//go:embed static/**
//go:embed static/*
//go:embed database/migrations/*
var embeddedFiles embed.FS

func main() {
	app := core.NewApp(embeddedFiles)
	app.Version = &VERSION
	cmd.ExecuteRootCommand(app)
}
