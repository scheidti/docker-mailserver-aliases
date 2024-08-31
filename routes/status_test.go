package routes

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDockerClient struct {
	mock.Mock
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)
	if args.Get(0) == nil {
		return []types.Container{}, args.Error(1)
	}
	return args.Get(0).([]types.Container), args.Error(1)
}

func (m *MockDockerClient) ContainerExecCreate(ctx context.Context, container string, config container.ExecOptions) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

func (m *MockDockerClient) ContainerExecAttach(ctx context.Context, execID string, config container.ExecAttachOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

func (m *MockDockerClient) Close() error {
	return nil
}

func TestStatusGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Docker client returns an error on ContainerList error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerList", mock.Anything, mock.Anything).Return(nil, errors.New("docker error"))

		router := gin.Default()
		router.GET("/v1/status", func(c *gin.Context) {
			checkIfContainerIsRunning(c, mockClient)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/status", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "docker error"}`, w.Body.String())
	})

	t.Run("No matching Docker containers running", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{
			{Image: "test/some-other-image"},
		}, nil)

		router := gin.Default()
		router.GET("/v1/status", func(c *gin.Context) {
			checkIfContainerIsRunning(c, mockClient)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/status", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"running": false}`, w.Body.String())
	})

	t.Run("Matching Docker container is running", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{
			{Image: "mailserver/docker-mailserver"},
		}, nil)

		router := gin.Default()
		router.GET("/v1/status", func(c *gin.Context) {
			checkIfContainerIsRunning(c, mockClient)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/status", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"running": true}`, w.Body.String())
	})

	t.Run("Matching Docker container from GitHub Container Registry is running", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{
			{Image: "ghcr.io/docker-mailserver/docker-mailserver"},
		}, nil)

		router := gin.Default()
		router.GET("/v1/status", func(c *gin.Context) {
			checkIfContainerIsRunning(c, mockClient)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/status", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"running": true}`, w.Body.String())
	})
}
