package main

import (
	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
)

func main() {

	app := core.NewApp()

	base_group := app.Gin.Group("/api/v1")

	// routes
	app.ResourceRoutes("projects", base_group, *handlers.NewProjectsHandler())
	app.ResourceRoutes("collections", base_group, *handlers.NewCollectionsHandler())
	app.ResourceRoutes("collections/:collection_id/fields", base_group, *handlers.NewFieldsHandler())
	app.ResourceRoutes("collections/:collection_id/records", base_group, *handlers.NewRecordsHandler())

	PORT := app.BaseUtil.SafeEnvGet("PORT", "3000")
	app.Gin.Run(":" + PORT)

}
