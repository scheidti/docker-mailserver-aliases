package routes

import (
	"bytes"
	"context"
	"io"
	"net/mail"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/scheidti/docker-mailserver-aliases/models"
)

// AliasesGetHandler godoc
//
//	@Summary	List of all available email aliases
//	@Schemes
//	@Description	Gets a list of all available email aliases from the Docker Mailserver container
//	@Tags			Aliases
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.AliasListResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/v1/aliases [get]
func AliasesGetHandler(c *gin.Context) {
	cli, err := getDockerClient()
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer cli.Close()

	container, err := getMailserverContainer(cli)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	aliases, err := getAliases(cli, container.ID)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, aliases)
}

func getAliases(cli DockerClient, containerName string) (models.AliasListResponse, error) {
	ctx := context.Background()

	execConfig := container.ExecOptions{
		Cmd:          []string{"setup", "alias", "list"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execId, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return models.AliasListResponse{}, err
	}

	resp, err := cli.ContainerExecAttach(ctx, execId.ID, container.ExecStartOptions{})
	if err != nil {
		return models.AliasListResponse{}, err
	}
	defer resp.Close()

	var outBuf, _ bytes.Buffer
	_, err = io.Copy(&outBuf, resp.Reader)
	if err != nil {
		return models.AliasListResponse{}, err
	}

	return parseAliasCommandResult(outBuf.String()), nil
}

func parseAliasCommandResult(commandResult string) models.AliasListResponse {
	lines := strings.Split(commandResult, "\n")
	splitRegex := regexp.MustCompile(`.*\* *([^\s]+) +([^\s]+)`)
	result := models.AliasListResponse{Aliases: make([]models.AliasResponse, 0)}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := splitRegex.FindStringSubmatch(line)

		if len(matches) == 3 {
			alias := matches[1]
			email := matches[2]
			_, aliasErr := mail.ParseAddress(alias)
			_, emailErr := mail.ParseAddress(email)
			if aliasErr != nil || emailErr != nil {
				continue
			}
			result.Aliases = append(result.Aliases, models.AliasResponse{Alias: alias, Email: email})
		}
	}

	return result
}
