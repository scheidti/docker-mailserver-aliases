package routes

import (
	"bytes"
	"context"
	"errors"
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

// AliasesPostHandler godoc
//
//	@Summary	Add a new email alias
//	@Schemes
//	@Description	Adds a new email alias to the Docker Mailserver container
//	@Tags			Aliases
//	@Accept			json
//	@Produce		json
//	@Param			alias	body		models.AliasResponse	true	"Alias to add"
//	@Success		201		{object}	models.AliasResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Router			/v1/aliases [post]
func AliasesPostHandler(c *gin.Context) {
	var newAlias models.AliasResponse

	if err := c.ShouldBindJSON(&newAlias); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

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

	_, err = checkIfAliasExists(cli, container.ID, newAlias.Alias)
	if err == nil {
		c.JSON(500, models.ErrorResponse{Error: "Alias already exists"})
		return
	}

	emailExists, err := checkIfEmailExists(cli, container.ID, newAlias.Email)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	if !emailExists {
		c.JSON(400, models.ErrorResponse{Error: "Email does not exist"})
		return
	}

	_, err = mail.ParseAddress(newAlias.Alias)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: "Invalid alias"})
		return
	}

	err = addAlias(cli, container.ID, newAlias)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(201, newAlias)
}

// AliasesDeleteHandler godoc
//
//	@Summary	Delete an email alias
//	@Schemes
//	@Description	Deletes an email alias from the Docker Mailserver container
//	@Tags			Aliases
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		500		{object}	models.ErrorResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Param			alias	path		string	true	"Alias to delete"
//	@Router			/v1/aliases/{alias} [delete]
func AliasesDeleteHandler(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(400, models.ErrorResponse{Error: "Alias must be provided"})
		return
	}

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

	existingAlias, err := checkIfAliasExists(cli, container.ID, alias)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = deleteAlias(cli, container.ID, existingAlias)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(204)
}

func checkIfAliasExists(cli DockerClient, containerName string, alias string) (models.AliasResponse, error) {
	aliases, err := getAliases(cli, containerName)
	if err != nil {
		return models.AliasResponse{}, err
	}

	for _, a := range aliases.Aliases {
		if a.Alias == alias {
			return a, nil
		}
	}

	return models.AliasResponse{}, errors.New("alias not found")
}

func checkIfEmailExists(cli DockerClient, containerName string, email string) (bool, error) {
	emails, err := getEmails(cli, containerName)
	if err != nil {
		return false, err
	}

	for _, e := range emails {
		if e == email {
			return true, nil
		}
	}

	return false, nil
}

func deleteAlias(cli DockerClient, containerName string, alias models.AliasResponse) error {
	ctx := context.Background()

	execConfig := container.ExecOptions{
		Cmd:          []string{"setup", "alias", "del", alias.Alias, alias.Email},
		AttachStdout: true,
		AttachStderr: true,
	}

	execId, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return err
	}

	_, err = cli.ContainerExecAttach(ctx, execId.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func addAlias(cli DockerClient, containerName string, alias models.AliasResponse) error {
	ctx := context.Background()

	execConfig := container.ExecOptions{
		Cmd:          []string{"setup", "alias", "add", alias.Alias, alias.Email},
		AttachStdout: true,
		AttachStderr: true,
	}

	execId, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return err
	}

	_, err = cli.ContainerExecAttach(ctx, execId.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}

	return nil
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
