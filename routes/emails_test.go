package routes

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestEmailsGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("getEmails should return a list of email addresses", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)

		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(bytes.NewBufferString("* name@developer.de ( 969K / ~ ) [0%] [ aliases -> postmaster@mail.de ]"))),
			Conn:   mockHijackedResponseConn,
		}, nil)

		emails, err := getEmails(mockClient, "containerId")
		assert.NoError(t, err)
		assert.Equal(t, []string{"name@developer.de"}, emails)
	})

	t.Run("getEmails should handle ContainerExecCreate error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{}, errors.New("exec create error"))

		emails, err := getEmails(mockClient, "containerId")
		assert.Error(t, err)
		assert.Nil(t, emails)
	})

	t.Run("getEmails should handle ContainerExecAttach error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{}, errors.New("exec attach error"))

		emails, err := getEmails(mockClient, "containerId")
		assert.Error(t, err)
		assert.Nil(t, emails)
	})

	t.Run("getEmails should handle io.Copy error", func(t *testing.T) {
		mockClient := new(MockDockerClient)
		mockHijackedResponseConn := new(MockHijackedResponseConn)
		mockHijackedResponseConn.On("Close").Return(nil)
		mockClient.On("ContainerExecCreate", mock.Anything, mock.Anything, mock.Anything).Return(types.IDResponse{ID: "execId"}, nil)
		mockClient.On("ContainerExecAttach", mock.Anything, mock.Anything, mock.Anything).Return(types.HijackedResponse{
			Reader: bufio.NewReader(io.NopCloser(&errorReader{})),
			Conn:   mockHijackedResponseConn,
		}, nil)

		emails, err := getEmails(mockClient, "containerId")
		assert.Error(t, err)
		assert.Nil(t, emails)
	})

	t.Run("parseEmailCommandResult should parse email addresses", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected []string
		}{
			{
				name:     "Valid emails",
				input:    "* test@example.com (alias)\n* user@domain.com (alias)",
				expected: []string{"test@example.com", "user@domain.com"},
			},
			{
				name:     "Invalid emails",
				input:    "* invalid-email (alias)\n* another-invalid-email (alias)",
				expected: []string{},
			},
			{
				name:     "Mixed valid and invalid emails",
				input:    "* valid@example.com (alias)\n* invalid-email (alias)\n* another@valid.com (alias)",
				expected: []string{"valid@example.com", "another@valid.com"},
			},
			{
				name:     "Empty input",
				input:    "",
				expected: []string{},
			},
			{
				name: "Complex input with multiple aliases",
				input: `* name@developer.de ( 969K / ~ ) [0%]
		[ aliases -> postmaster@mail.de ]
	
	* admin@website.net ( 2.5M / ~ ) [0%]
		[ aliases -> postmaster@website.net, webmaster@website.net ]
	
	* name@company.tech ( 16M / ~ ) [0%]
		[ aliases -> postmaster@company.tech, shop@company.tech, admin@company.tech ]`,
				expected: []string{
					"name@developer.de",
					"admin@website.net",
					"name@company.tech",
				},
			},
			{
				name:     "Single valid email",
				input:    "* single@valid.com (alias)",
				expected: []string{"single@valid.com"},
			},
			{
				name:     "Single invalid email",
				input:    "* invalid-email (alias)",
				expected: []string{},
			},
			{
				name:     "Emails with special characters",
				input:    "* special+chars@example.com (alias)\n* another.special@domain.com (alias)",
				expected: []string{"special+chars@example.com", "another.special@domain.com"},
			},
			{
				name:     "Emails with different casing",
				input:    "* MixedCase@Example.com (alias)\n* another@Domain.COM (alias)",
				expected: []string{"MixedCase@Example.com", "another@Domain.COM"},
			},
			{
				name:     "Emails with leading or trailing spaces",
				input:    "*  leading@space.com (alias)\n* trailing@space.com  (alias)",
				expected: []string{"leading@space.com", "trailing@space.com"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := parseEmailCommandResult(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}
