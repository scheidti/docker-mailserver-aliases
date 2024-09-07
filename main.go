package main

import (
	"embed"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scheidti/docker-mailserver-aliases/docs"
	"github.com/scheidti/docker-mailserver-aliases/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Docker Mailserver Aliases API
//	@version		1.0
//	@description	API for managing aliases in a Docker Mailserver container

//	@contact.name	Christian Scheid
//	@contact.url	https://github.com/scheidti/docker-mailserver-aliases
//	@contact.email	admin@scheid.tech

//	@license.name	MIT
//	@license.url	https://github.com/scheidti/docker-mailserver-aliases?tab=MIT-1-ov-file#readme

//	@host	localhost:8080

//go:embed frontend/dist/*
var frontend embed.FS

func main() {
	engine := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	api := engine.Group("/v1")
	{
		api.GET("/status", routes.StatusGetHandler)
		api.GET("/emails", routes.EmailsGetHandler)
		api.GET("/aliases", routes.AliasesGetHandler)
		api.POST("/aliases", routes.AliasesPostHandler)
		api.DELETE("/aliases/:alias", routes.AliasesDeleteHandler)
	}

	addr := os.Getenv("GIN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	if gin.Mode() != gin.ReleaseMode {
		engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	engine.NoRoute(serveFrontend)
	engine.Run(addr)
}

func serveFrontend(c *gin.Context) {
	fileServer := http.StripPrefix("/", http.FileServer(http.FS(frontend)))
	c.Request.URL.Path = "/frontend/dist" + c.Request.URL.Path
	fileServer.ServeHTTP(c.Writer, c.Request)
}
