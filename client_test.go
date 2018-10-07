package ts3

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second))
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.Error(t, c.Close())
	}()

	_, err = c.Exec("version")
	assert.NoError(t, err)

	_, err = c.ExecCmd(NewCmd("version"))
	assert.NoError(t, err)

	_, err = c.ExecCmd(NewCmd("invalid"))
	assert.Error(t, err)

	_, err = c.ExecCmd(NewCmd("disconnect"))
	assert.Error(t, err)
}

func TestClientNilOption(t *testing.T) {
	_, err := NewClient("", nil)
	if !assert.Error(t, err) {
		return
	}

	assert.Equal(t, ErrNilOption, err)
}

func TestClientOptionError(t *testing.T) {
	errBadOption := errors.New("bad option")
	_, err := NewClient("", func(c *Client) error { return errBadOption })
	if !assert.Error(t, err) {
		return
	}

	assert.Equal(t, errBadOption, err)
}

func TestClientDisconnect(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, c.Close())

	_, err = c.Exec("version")
	assert.Error(t, err)
}

func TestClientWriteFail(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second))
	if !assert.NoError(t, err) {
		return
	}
	assert.NoError(t, c.conn.(*net.TCPConn).CloseWrite())

	_, err = c.Exec("version")
	assert.Error(t, err)
}

func TestClientDialFail(t *testing.T) {
	c, err := NewClient("127.0.0.1", Timeout(time.Nanosecond))
	if assert.Error(t, err) {
		return
	}

	// Should never get here
	assert.NoError(t, c.Close())
}

func TestClientTimeout(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Millisecond*100))
	if !assert.NoError(t, err) {
		return
	}

	// Not receiving a response must cause a timeout
	_, err = c.Exec(" ")
	assert.Error(t, err)
}

func TestClientDeadline(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Millisecond*100))
	if !assert.NoError(t, err) {
		return
	}

	_, err = c.Exec("version")
	assert.NoError(t, err)

	// Inactivity must not cause a timeout
	time.Sleep(c.timeout * 2)

	_, err = c.Exec("version")
	assert.NoError(t, err)
}

func TestClientNoHeader(t *testing.T) {
	s := newServerStopped(t)
	if s == nil {
		return
	}
	s.noHeader = true
	s.Start()
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Millisecond*100))
	if assert.Error(t, err) {
		return
	}

	// Should never get here
	assert.NoError(t, c.Close())
}

func TestClientNoBanner(t *testing.T) {
	s := newServerStopped(t)
	if s == nil {
		return
	}
	s.noBanner = true
	s.Start()
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Millisecond*100))
	if assert.Error(t, err) {
		return
	}

	// Should never get here
	assert.NoError(t, c.Close())
}

func TestClientFailConn(t *testing.T) {
	s := newServerStopped(t)
	if s == nil {
		return
	}
	s.failConn = true
	s.Start()
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second))
	if assert.Error(t, err) {
		return
	}

	// Should never get here
	assert.NoError(t, c.Close())
}

func TestClientBadHeader(t *testing.T) {
	s := newServerStopped(t)
	if s == nil {
		return
	}
	s.badHeader = true
	s.Start()
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second))
	if assert.Error(t, err) {
		return
	}

	// Should never get here
	assert.NoError(t, c.Close())
}
