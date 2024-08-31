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

// EmailsGetHandler godoc
//
//	@Summary	List of all available email addresses
//	@Schemes
//	@Description	Gets a list of all available email addresses from the Docker Mailserver container
//	@Tags			E-Mails
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.EmailListResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/v1/emails [get]
func EmailsGetHandler(c *gin.Context) {
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

	emails, err := getEmails(cli, container.ID)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, models.EmailListResponse{Emails: emails})
}

func getEmails(cli DockerClient, containerName string) ([]string, error) {
	ctx := context.Background()

	execConfig := container.ExecOptions{
		Cmd:          []string{"setup", "email", "list"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execId, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerExecAttach(context.Background(), execId.ID, container.ExecStartOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var outBuf, _ bytes.Buffer
	_, err = io.Copy(&outBuf, resp.Reader)
	if err != nil {
		return nil, err
	}

	return parseEmailCommandResult(outBuf.String()), nil
}

func parseEmailCommandResult(commandResult string) []string {
	emailRegex := regexp.MustCompile(`\*\s*(.*?)\s*\(`)
	lines := strings.Split(commandResult, "\n")
	result := make([]string, 0)

	for _, line := range lines {
		matches := emailRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			_, err := mail.ParseAddress(matches[1])
			if err != nil {
				continue
			}
			result = append(result, matches[1])
		}
	}

	return result
}
