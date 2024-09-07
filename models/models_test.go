package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerImageConstant(t *testing.T) {
	assert.Equal(t, "mailserver/docker-mailserver", DockerImage, "DockerImage should match the expected value")
}

func TestStatusResponseMarshalling(t *testing.T) {
	original := StatusResponse{Running: true}
	data, err := json.Marshal(original)
	assert.NoError(t, err, "Marshalling StatusResponse should not return an error")

	expectedJSON := `{"running":true}`
	assert.JSONEq(t, expectedJSON, string(data), "Marshalled JSON should match expected")

	var unmarshalled StatusResponse
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err, "Unmarshalling StatusResponse should not return an error")
	assert.Equal(t, original, unmarshalled, "Unmarshalled StatusResponse should match original")
}

func TestErrorResponseMarshalling(t *testing.T) {
	original := ErrorResponse{Error: "Something went wrong"}
	data, err := json.Marshal(original)
	assert.NoError(t, err, "Marshalling ErrorResponse should not return an error")

	expectedJSON := `{"error":"Something went wrong"}`
	assert.JSONEq(t, expectedJSON, string(data), "Marshalled JSON should match expected")

	var unmarshalled ErrorResponse
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err, "Unmarshalling ErrorResponse should not return an error")
	assert.Equal(t, original, unmarshalled, "Unmarshalled ErrorResponse should match original")
}

func TestEmailListResponseMarshalling(t *testing.T) {
	original := EmailListResponse{Emails: []string{"user1@example.com", "user2@example.com"}}
	data, err := json.Marshal(original)
	assert.NoError(t, err, "Marshalling EmailListResponse should not return an error")

	expectedJSON := `{"emails":["user1@example.com","user2@example.com"]}`
	assert.JSONEq(t, expectedJSON, string(data), "Marshalled JSON should match expected")

	var unmarshalled EmailListResponse
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err, "Unmarshalling EmailListResponse should not return an error")
	assert.Equal(t, original, unmarshalled, "Unmarshalled EmailListResponse should match original")
}

func TestAliasResponseMarshalling(t *testing.T) {
	original := AliasResponse{Alias: "alias@example.com", Email: "user@example.com"}
	data, err := json.Marshal(original)
	assert.NoError(t, err, "Marshalling AliasResponse should not return an error")

	expectedJSON := `{"alias":"alias@example.com","email":"user@example.com"}`
	assert.JSONEq(t, expectedJSON, string(data), "Marshalled JSON should match expected")

	var unmarshalled AliasResponse
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err, "Unmarshalling AliasResponse should not return an error")
	assert.Equal(t, original, unmarshalled, "Unmarshalled AliasResponse should match original")
}

func TestAliasListResponseMarshalling(t *testing.T) {
	original := AliasListResponse{
		Aliases: []AliasResponse{
			{Alias: "alias1@example.com", Email: "user1@example.com"},
			{Alias: "alias2@example.com", Email: "user2@example.com"},
		},
	}
	data, err := json.Marshal(original)
	assert.NoError(t, err, "Marshalling AliasListResponse should not return an error")

	expectedJSON := `{"aliases":[{"alias":"alias1@example.com","email":"user1@example.com"},{"alias":"alias2@example.com","email":"user2@example.com"}]}`
	assert.JSONEq(t, expectedJSON, string(data), "Marshalled JSON should match expected")

	var unmarshalled AliasListResponse
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err, "Unmarshalling AliasListResponse should not return an error")
	assert.Equal(t, original, unmarshalled, "Unmarshalled AliasListResponse should match original")
}
