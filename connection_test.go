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
