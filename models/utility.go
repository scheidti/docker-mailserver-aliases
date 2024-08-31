package models

const DockerImage = "mailserver/docker-mailserver"

type StatusResponse struct {
	Running bool `json:"running"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmailListResponse struct {
	Emails []string `json:"emails"`
}
