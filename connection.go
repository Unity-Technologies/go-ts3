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

// legacyConnection is an insecure TCP connection.
type legacyConnection struct {
	net.Conn
}

// Connect connects to the address with the given timeout.
func (c *legacyConnection) Connect(addr string, timeout time.Duration) error {
	addr, err := verifyAddr(addr, DefaultPort)
	if err != nil {
		return err
	}

	c.Conn, err = net.DialTimeout("tcp", addr, timeout)
	return err
}

// sshConnection is an SSH connection with open SSH channel and attached shell.
type sshConnection struct {
	net.Conn
	config  *ssh.ClientConfig
	channel ssh.Channel
}

// Connect connects to the address with the given timeout and opens a new SSH channel with attached shell.
func (c *sshConnection) Connect(addr string, timeout time.Duration) error {
	addr, err := verifyAddr(addr, DefaultSSHPort)
	if err != nil {
		return err
	}

	if c.Conn, err = net.DialTimeout("tcp", addr, timeout); err != nil {
		return err
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(c.Conn, addr, c.config)
	if err != nil {
		return err
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
		return err
	}
	go ssh.DiscardRequests(reqs)

	ok, err := c.channel.SendRequest("shell", true, nil)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("ssh connection: could not open shell")
	}

	return nil
}

// Read implements io.Reader.
func (c *sshConnection) Read(p []byte) (n int, err error) {
	return c.channel.Read(p)
}

// Write implements io.Writer.
func (c *sshConnection) Write(p []byte) (n int, err error) {
	return c.channel.Write(p)
}

// Close implements io.Closer.
func (c *sshConnection) Close() error {
	var err error
	if err2 := c.channel.Close(); err2 != nil && !errors.Is(err2, io.EOF) {
		err = err2
	}
	if err2 := c.Conn.Close(); err2 != nil && !errors.Is(err2, net.ErrClosed) {
		err = err2
	}
	return err
}

// verifyAddr checks if addr is formatted correctly. If valid it returns addr.
// If the address does not include a port, defaultPort is added.
// A literal IPv6 must be enclosed in square brackets e.g. "[::1]"
func verifyAddr(addr string, defaultPort int) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok && addrError.Err == "missing port in address" {
			return net.JoinHostPort(strings.Trim(addr, "[]"), strconv.Itoa(defaultPort)), nil
		}
		return "", err
	}
	return net.JoinHostPort(host, port), nil
}
