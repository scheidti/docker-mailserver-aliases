package routes

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/scheidti/docker-mailserver-aliases/models"
)

// statusGET godoc
//
//	@Summary	Checks Mailserver Docker container
//	@Schemes
//	@Description	Checks if the Docker Mailserver Docker container is running
//	@Tags			Utility
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.StatusResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/v1/status [get]
func StatusGET(c *gin.Context) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	for _, container := range containers {
		if strings.Contains(container.Image, models.DockerImage) {
			c.JSON(200, models.StatusResponse{Running: true})
			return
		}
	}

	c.JSON(200, models.StatusResponse{Running: false})
}
