package ts3

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	// DefaultPort is the default TeamSpeak 3 ServerQuery port.
	DefaultPort = 10011

	// DefaultSSHPort is the default TeamSpeak 3 ServerQuery SSH port.
	DefaultSSHPort = 10022
)

// writeTimeoutConn is a net.Conn that sets the write timeout on every call to Write().
type writeTimeoutConn struct {
	net.Conn
	timeout time.Duration
}

func (c *writeTimeoutConn) Write(p []byte) (n int, err error) {
	if err = c.Conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, fmt.Errorf("writeTimeoutConn: SetWriteDeadline: %w", err)
	}
	if n, err = c.Conn.Write(p); err != nil {
		return n, fmt.Errorf("writeTimeoutConn: write: %w", err)
	}
	return n, nil
}

// legacyConnection is an insecure TCP connection.
type legacyConnection struct {
	net.Conn
}

// Connect connects to the address with the given timeout.
// The timeout is used as dial and write timeout.
func (c *legacyConnection) Connect(addr string, timeout time.Duration) error {
	addr, err := verifyAddr(addr, DefaultPort)
	if err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("legacy connection: dial: %w", err)
	}

	c.Conn = &writeTimeoutConn{
		Conn:    conn,
		timeout: timeout,
	}

	return nil
}

// sshConnection is an SSH connection with open SSH channel and attached shell.
type sshConnection struct {
	net.Conn
	config  *ssh.ClientConfig
	channel ssh.Channel
}

// Connect connects to the address with the given timeout and opens a new SSH channel with attached shell.
// The timeout is used as dial and write timeout.
func (c *sshConnection) Connect(addr string, timeout time.Duration) error {
	addr, err := verifyAddr(addr, DefaultSSHPort)
	if err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("ssh connection: dial: %w", err)
	}

	c.Conn = &writeTimeoutConn{
		Conn:    conn,
		timeout: timeout,
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(c.Conn, addr, c.config)
	if err != nil {
		return fmt.Errorf("ssh connecion: ssh client conn: %w", err)
	}
	go ssh.DiscardRequests(reqs)

	// Reject all channel requests.
	go func(newChannel <-chan ssh.NewChannel) {
		for channel := range newChannel {
			channel.Reject(ssh.Prohibited, ssh.Prohibited.String()) //nolint: errcheck
		}
	}(chans)

	c.channel, reqs, err = clientConn.OpenChannel("session", nil)
	if err != nil {
		return fmt.Errorf("ssh connection: session: %w", err)
	}
	go ssh.DiscardRequests(reqs)

	ok, err := c.channel.SendRequest("shell", true, nil)
	if err != nil {
		return fmt.Errorf("ssh connection: shell: %w", err)
	}
	if !ok {
		return fmt.Errorf("ssh connection: could not open shell")
	}

	return nil
}

// Read implements io.Reader.
func (c *sshConnection) Read(p []byte) (n int, err error) {
	if n, err = c.channel.Read(p); err != nil {
		return n, fmt.Errorf("ssh connection: read: %w", err)
	}
	return n, nil
}

// Write implements io.Writer.
func (c *sshConnection) Write(p []byte) (n int, err error) {
	if n, err = c.channel.Write(p); err != nil {
		return n, fmt.Errorf("ssh connection: write: %w", err)
	}
	return n, nil
}

// Close implements io.Closer.
func (c *sshConnection) Close() error {
	var err error
	// In both cases we ignore errors which don't have any value.
	if err2 := c.channel.Close(); err2 != nil &&
		!errors.Is(err2, io.EOF) &&
		!strings.HasSuffix(err2.Error(), "connection reset by peer") {
		err = err2
	}
	if err2 := c.Conn.Close(); err2 != nil &&
		err == nil &&
		!errors.Is(err2, net.ErrClosed) &&
		!strings.HasSuffix(err2.Error(), "connection reset by peer") {
		err = err2
	}
	return err
}

// verifyAddr checks if addr is formatted correctly. If valid it returns addr.
// If the address does not include a port, defaultPort is added.
// A literal IPv6 must be enclosed in square brackets e.g. "[::1]".
func verifyAddr(addr string, defaultPort int) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		var addrError *net.AddrError
		if ok := errors.As(err, &addrError); ok && addrError.Err == "missing port in address" {
			return net.JoinHostPort(strings.Trim(addr, "[]"), strconv.Itoa(defaultPort)), nil
		}
		return "", fmt.Errorf("verify address: %w", err)
	}
	return net.JoinHostPort(host, port), nil
}
