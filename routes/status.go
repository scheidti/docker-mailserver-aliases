package routes

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/scheidti/docker-mailserver-aliases/models"
)

type DockerClient interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerExecCreate(ctx context.Context, container string, config container.ExecOptions) (types.IDResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config container.ExecAttachOptions) (types.HijackedResponse, error)
	Close() error
}

// StatusGetHandler godoc
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
func StatusGetHandler(c *gin.Context) {
	cli, err := getDockerClient()
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer cli.Close()

	checkIfContainerIsRunning(c, cli)
}

func getDockerClient() (DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func checkIfContainerIsRunning(c *gin.Context, cli DockerClient) {
	ctx := context.Background()

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

func getMailserverContainer(cli DockerClient) (types.Container, error) {
	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return types.Container{}, err
	}

	for _, container := range containers {
		if strings.Contains(container.Image, models.DockerImage) {
			return container, nil
		}
	}

	return types.Container{}, errors.New("mailserver container not found")
}
