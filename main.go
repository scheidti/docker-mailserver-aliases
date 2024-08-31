package main

import (
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
//	@contact.email	christian@scheid.tech

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host	localhost:8080

func main() {
	engine := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	api := engine.Group("/v1")
	{
		api.GET("/status", routes.StatusGetHandler)
		api.GET("/emails", routes.EmailsGetHandler)
		api.GET("/aliases", routes.AliasesGetHandler)
	}

	addr := os.Getenv("GIN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.Run(addr)
}
