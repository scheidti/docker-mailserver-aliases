package routes

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/scheidti/docker-mailserver-aliases/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHijackedResponseConn struct {
	mock.Mock
}

func (e *MockHijackedResponseConn) Close() error {
	return nil
}

func (e *MockHijackedResponseConn) LocalAddr() net.Addr {
	return nil
}

func (e *MockHijackedResponseConn) RemoteAddr() net.Addr {
	return nil
}

func (e *MockHijackedResponseConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (e *MockHijackedResponseConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (e *MockHijackedResponseConn) SetDeadline(t time.Time) error {
	return nil
}

func (e *MockHijackedResponseConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (e *MockHijackedResponseConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestAliasGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("getAliases should return a list of aliases", func(t *testing.T) {
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(bytes.NewBufferString(`* postmaster@website.de admin@website.de
* alias2@website.de admin@website.de`))),
			Conn: mockHijackedResponseConn,
		}, nil)

		aliases, err := getAliases(mockClient, "containerId")
		assert.NoError(t, err)
		assert.Equal(t, models.AliasListResponse{
			Aliases: []models.AliasResponse{
				{Alias: "postmaster@website.de", Email: "admin@website.de"},
				{Alias: "alias2@website.de", Email: "admin@website.de"},
			},
		}, aliases)
	})

	t.Run("getAliases should handle ContainerExecCreate error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{}, errors.New("exec create error"))

		aliases, err := getAliases(mockClient, "containerId")
		assert.Error(t, err)
		assert.Empty(t, aliases.Aliases)
	})

	t.Run("getAliases should handle ContainerExecAttach error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{}, errors.New("exec attach error"))

		aliases, err := getAliases(mockClient, "containerId")
		assert.Error(t, err)
		assert.Empty(t, aliases.Aliases)
	})

	t.Run("getAliases should handle io.Copy error", func(t *testing.T) {
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(&errorReader{})),
			Conn:   mockHijackedResponseConn,
		}, nil)

		aliases, err := getAliases(mockClient, "containerId")
		assert.Error(t, err)
		assert.Empty(t, aliases.Aliases)
	})

	t.Run("parseAliasCommandResult should parse aliases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected models.AliasListResponse
		}{
			{
				name:  "single alias",
				input: "* postmaster@website.de admin@website.de",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{
						{Alias: "postmaster@website.de", Email: "admin@website.de"},
					},
				},
			},
			{
				name:  "multiple aliases",
				input: "* postmaster@v-developer.de cscheid@v-developer.de\n* postmaster@scheid.tech christian@scheid.tech",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{
						{Alias: "postmaster@v-developer.de", Email: "cscheid@v-developer.de"},
						{Alias: "postmaster@scheid.tech", Email: "christian@scheid.tech"},
					},
				},
			},
			{
				name:  "empty input",
				input: "",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{},
				},
			},
			{
				name:  "alias with extra spaces",
				input: " *  postmaster@scheidti.net    admin@scheidti.net ",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{
						{Alias: "postmaster@scheidti.net", Email: "admin@scheidti.net"},
					},
				},
			},
			{
				name:  "alias with mixed case",
				input: "* Postmaster@Example.com Admin@Example.com",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{
						{Alias: "Postmaster@Example.com", Email: "Admin@Example.com"},
					},
				},
			},
			{
				name:  "alias with multiple lines and mixed case",
				input: "* Postmaster@Example.com Admin@Example.com\n* amazon@Example.net User@Example.net",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{
						{Alias: "Postmaster@Example.com", Email: "Admin@Example.com"},
						{Alias: "amazon@Example.net", Email: "User@Example.net"},
					},
				},
			},
			{
				name:  "multiple aliases, varying domains",
				input: "* admin@domain1.com user@domain1.com\n* postmaster@domain2.com admin@domain2.com\n* support@domain3.com user@domain3.com",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{
						{Alias: "admin@domain1.com", Email: "user@domain1.com"},
						{Alias: "postmaster@domain2.com", Email: "admin@domain2.com"},
						{Alias: "support@domain3.com", Email: "user@domain3.com"},
					},
				},
			},
			{
				name:  "alias with no email",
				input: "* postmaster@website.de",
				expected: models.AliasListResponse{
					Aliases: []models.AliasResponse{},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := parseAliasCommandResult(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestAliasPostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("addAlias should add an alias", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{}, nil)

		err := addAlias(mockClient, "containerId", models.AliasResponse{Alias: "test@alias.de", Email: "user@mail.de"})
		assert.NoError(t, err)
	})

	t.Run("addAlias should handle ContainerExecCreate error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{}, errors.New("exec create error"))

		err := addAlias(mockClient, "containerId", models.AliasResponse{Alias: "test@alias.de", Email: "user@mail.de"})
		assert.Error(t, err)
	})

	t.Run("addAlias should handle ContainerExecAttach error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{}, errors.New("exec attach error"))

		err := addAlias(mockClient, "containerId", models.AliasResponse{Alias: "test@alias.de", Email: "user@mail.de"})
		assert.Error(t, err)
	})

	t.Run("checkIfEmailExists should return true if email exists", func(t *testing.T) {
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(bytes.NewBufferString(`* name@developer.de ( 969K / ~ ) [0%] [ aliases -> postmaster@mail.de ]`))),
			Conn:   mockHijackedResponseConn,
		}, nil)

		exists, err := checkIfEmailExists(mockClient, "containerId", "name@developer.de")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("checkIfEmailExists should return false if email does not exist", func(t *testing.T) {
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(bytes.NewBufferString(`* name@developer.de ( 969K / ~ ) [0%] [ aliases -> postmaster@mail.de ]`))),
			Conn:   mockHijackedResponseConn,
		}, nil)

		exists, err := checkIfEmailExists(mockClient, "containerId", "doesNotExist@developer.de")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("checkIfAliasExists should return the alias response if alias exists", func(t *testing.T) {
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(bytes.NewBufferString(`* postmaster@website.de admin@website.de
* alias2@website.de admin@website.de`))),
			Conn: mockHijackedResponseConn,
		}, nil)

		exists, err := checkIfAliasExists(mockClient, "containerId", "alias2@website.de")
		assert.NoError(t, err)
		assert.Equal(t, models.AliasResponse{Alias: "alias2@website.de", Email: "admin@website.de"}, exists)
	})

	t.Run("checkIfAliasExists should return an error if alias does not exist", func(t *testing.T) {
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(bytes.NewBufferString(`* postmaster@website.de admin@website.de
* alias2@website.de admin@website.de`))),
			Conn: mockHijackedResponseConn,
		}, nil)

		exists, err := checkIfAliasExists(mockClient, "containerId", "wrong@website.de")
		assert.Error(t, err)
		assert.Empty(t, exists.Alias)
	})

	t.Run("POST with invalid JSON should return 400", func(t *testing.T) {
		router := gin.Default()
		router.POST("/v1/aliases", func(c *gin.Context) {
			AliasesPostHandler(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/aliases", bytes.NewBufferString(`No JSON`))
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"error": "Invalid request body"}`, w.Body.String())
	})

	t.Run("POST with invalid alias should return 400", func(t *testing.T) {
		router := gin.Default()
		router.POST("/v1/aliases", func(c *gin.Context) {
			AliasesPostHandler(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/aliases", bytes.NewBufferString(`{"alias": "invalid"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"error": "Invalid alias"}`, w.Body.String())
	})
}

func TestAliasDeleteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("deleteAlias should delete an alias", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{}, nil)

		err := deleteAlias(mockClient, "containerId", models.AliasResponse{Alias: "alias@mail.de", Email: "user@mail.de"})
		assert.NoError(t, err)
	})

	t.Run("deleteAlias should handle ContainerExecCreate error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{}, errors.New("exec create error"))

		err := deleteAlias(mockClient, "containerId", models.AliasResponse{Alias: "alias@mail.de", Email: "user@mail.de"})
		assert.Error(t, err)
	})

	t.Run("deleteAlias should handle ContainerExecAttach error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{}, errors.New("exec attach error"))

		err := deleteAlias(mockClient, "containerId", models.AliasResponse{Alias: "alias@mail.de", Email: "user@mail.de"})
		assert.Error(t, err)
	})

}
