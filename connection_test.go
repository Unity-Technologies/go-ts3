package ts3

import (
	"bufio"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

var sshClientTestConfig = &ssh.ClientConfig{HostKeyCallback: ssh.InsecureIgnoreHostKey()} //nolint:gosec

func TestConnection(t *testing.T) {
	testCases := map[string]struct {
		conn      Connection
		newServer func(*testing.T) *server
	}{
		"legacyConnection": {
			conn:      new(legacyConnection),
			newServer: newServer,
		},
		"sshConnection": {
			conn: &sshConnection{config: sshClientTestConfig},
			newServer: func(t *testing.T) *server {
				s := newServer(t)
				s.useSSH = true
				return s
			},
		},
	}

	for description, tc := range testCases {
		t.Run(description, func(t *testing.T) {
			s := tc.newServer(t)
			require.NotNil(t, s)

			assert.NoError(t, tc.conn.Connect(s.Addr, time.Second))

			line, _, err := bufio.NewReader(tc.conn).ReadLine()
			assert.NoError(t, err)
			assert.Equal(t, "TS3", string(line))

			assert.NoError(t, s.Close())
			assert.NoError(t, tc.conn.Close())
		})
	}
}

func TestSplitHostAndPort(t *testing.T) {
	testCases := map[string]struct {
		addr         string
		expectError  bool
		expectedHost string
		expectedPort string
	}{
		"hostname without port": {
			addr:         "localhost",
			expectedHost: "localhost",
			expectedPort: "1234",
		},
		"hostname with port": {
			addr:         "localhost:1337",
			expectedHost: "localhost",
			expectedPort: "1337",
		},
		"IPv4 without port": {
			addr:         "127.0.0.1",
			expectedHost: "127.0.0.1",
			expectedPort: "1234",
		},
		"IPv4 with port": {
			addr:         "127.0.0.1:1337",
			expectedHost: "127.0.0.1",
			expectedPort: "1337",
		},
		"IPv6 without port": {
			addr:         "[::1]",
			expectedHost: "::1",
			expectedPort: "1234",
		},
		"IPv6 with port": {
			addr:         "[::1]:1337",
			expectedHost: "::1",
			expectedPort: "1337",
		},
		"IPv6 without brackets": {
			addr:        "::1",
			expectError: true,
		},
		"empty addr": {
			addr:         "",
			expectedHost: "",
			expectedPort: "1234",
		},
	}

	for description, tc := range testCases {
		t.Run(description, func(t *testing.T) {
			host, port, err := splitHostAndPort(tc.addr, 1234)
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tc.expectedHost, host)
				assert.Equal(t, tc.expectedPort, port)
				assert.NoError(t, err)
			}
		})
	}
}
