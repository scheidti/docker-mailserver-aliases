package models

import "os"

func GetDockerImage() string {
	if image := os.Getenv("DOCKER_MAILSERVER_IMAGE"); image != "" {
		return image
	}
	return "mailserver/docker-mailserver"
}

type StatusResponse struct {
	Running bool `json:"running"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmailListResponse struct {
	Emails []string `json:"emails"`
}

type AliasListResponse struct {
	Aliases []AliasResponse `json:"aliases"`
}

type AliasResponse struct {
	Alias string `json:"alias"`
	Email string `json:"email"`
}
