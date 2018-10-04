package ts3

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultPort is the default TeamSpeak 3 ServerQuery port.
	DefaultPort = 10011

	// MaxParseTokenSize is the maximum buffer size used to parse the
	// server responses.
	// It's relatively large to enable us to deal with the typical responses
	// to commands such as serversnapshotcreate.
	MaxParseTokenSize = 10 << 20

	// connectHeader is the header used as the prefix to responses.
	connectHeader = "TS3"

	// startBufSize is the initial size of allocation for the parse buffer.
	startBufSize = 4096

	// defaultNotificationBufSize is the initial notification buffer size.
	defaultNotificationBufSize = 5
)

var (
	respTrailerRe = regexp.MustCompile(`^error id=(\d+) msg=([^ ]+)(.*)`)

	// DefaultTimeout is the default read / write / dial timeout for Clients.
	DefaultTimeout = 10 * time.Second

	// keepaliveInterval is the interval in which keepalive data is sent.
	keepaliveInterval = 200 * time.Second
)

// Client is a TeamSpeak 3 ServerQuery client.
type Client struct {
	conn       net.Conn
	timeout    time.Duration
	scanner    *bufio.Scanner
	buf        []byte
	maxBufSize int
	mutex      sync.Mutex
	notify     chan Notification
	err        chan error
	disconnect chan bool
	res        []string

	Server *ServerMethods
}

// Timeout sets read / write / dial timeout for a TeamSpeak 3 Client.
func Timeout(timeout time.Duration) func(*Client) error {
	return func(c *Client) error {
		c.timeout = timeout
		return nil
	}
}

// Keepalive keeps the connection open.
func Keepalive() func(*Client) error {
	return func(c *Client) error {
		go func(c *Client) {
			for c.IsConnected() {
				time.Sleep(keepaliveInterval)

				if err := c.setDeadline(); err != nil {
					break
				}

				if _, err := c.conn.Write([]byte(" \n")); err != nil {
					break
				}

				if err := c.clearDeadline(); err != nil {
					break
				}
			}

			c.setDisconnected()
		}(c)
		return nil
	}
}

// NotificationBuffer sets the notification buffer size.
func NotificationBuffer(size int) func(*Client) error {
	return func(c *Client) error {
		c.notify = make(chan Notification, size)
		return nil
	}
}

// Buffer sets the initial buffer used to parse responses from
// the server and the maximum size of buffer that may be allocated.
// The maximum parsable token size is the larger of max and cap(buf).
// If max <= cap(buf), scanning will use this buffer only and do no
// allocation.
//
// By default, parsing uses an internal buffer and sets the maximum
// token size to MaxParseTokenSize.
func Buffer(buf []byte, max int) func(*Client) error {
	return func(c *Client) error {
		c.buf = buf
		c.maxBufSize = max
		return nil
	}
}

// NewClient returns a new TeamSpeak 3 client connected to addr.
func NewClient(addr string, options ...func(c *Client) error) (*Client, error) {
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%v:%v", addr, DefaultPort)
	}

	c := &Client{
		timeout:    DefaultTimeout,
		buf:        make([]byte, startBufSize),
		maxBufSize: MaxParseTokenSize,
		notify:     make(chan Notification, defaultNotificationBufSize),
		err:        make(chan error),
		disconnect: make(chan bool),
	}
	for _, f := range options {
		if f == nil {
			return nil, ErrNilOption
		}
		if err := f(c); err != nil {
			return nil, err
		}
	}

	// Wire up command groups
	c.Server = &ServerMethods{Client: c}

	var err error
	if c.conn, err = net.DialTimeout("tcp", addr, c.timeout); err != nil {
		return nil, err
	}

	c.scanner = bufio.NewScanner(bufio.NewReader(c.conn))
	c.scanner.Buffer(c.buf, c.maxBufSize)
	c.scanner.Split(ScanLines)

	if err := c.setDeadline(); err != nil {
		return nil, err
	}

	// Read the connection header
	if !c.scanner.Scan() {
		return nil, c.scanErr()
	}

	if l := c.scanner.Text(); l != connectHeader {
		return nil, fmt.Errorf("invalid connection header %q", l)
	}

	// Slurp the banner
	if !c.scanner.Scan() {
		return nil, c.scanErr()
	}

	if err := c.clearDeadline(); err != nil {
		return nil, err
	}

	// Handle incoming lines
	go c.messageHandler()

	return c, nil
}

// messageHandler scans incoming lines and handles them accordingly.
func (c *Client) messageHandler() {
	for c.IsConnected() {
		if c.scanner.Scan() {
			line := c.scanner.Text()
			if line == "error id=0 msg=ok" {
				c.err <- nil
			} else if matches := respTrailerRe.FindStringSubmatch(line); len(matches) == 4 {
				c.err <- NewError(matches)
			} else if strings.Index(line, "notify") == 0 {
				if n, err := decodeNotification(line); err == nil {
					// non-blocking write
					select {
					case c.notify <- n:
					default:
					}
				}
			} else {
				c.res = append(c.res, line)
			}
		} else if c.scanner.Err() == nil {
			c.err <- c.scanErr()
		} else {
			break
		}
	}

	c.setDisconnected()
	c.err <- c.scanErr()
}

// setDisconnected marks the clients as disconnected.
func (c *Client) setDisconnected() {
	select {
	case c.disconnect <- true:
	default:
	}
}

// setDeadline updates the deadline on the connection based on the clients configured timeout.
func (c *Client) setDeadline() error {
	return c.conn.SetDeadline(time.Now().Add(c.timeout))
}

// clearDeadline clears the deadline on the connection.
func (c *Client) clearDeadline() error {
	return c.conn.SetDeadline(time.Time{})
}

// Exec executes cmd on the server and returns the response.
func (c *Client) Exec(cmd string) ([]string, error) {
	return c.ExecCmd(NewCmd(cmd))
}

// ExecCmd executes cmd on the server and returns the response.
func (c *Client) ExecCmd(cmd *Cmd) ([]string, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.res = nil

	if err := c.setDeadline(); err != nil {
		return nil, err
	}

	if _, err := c.conn.Write([]byte(cmd.String())); err != nil {
		return nil, err
	}

	if err := c.setDeadline(); err != nil {
		return nil, err
	}

	select {
	case err := <-c.err:
		if err != nil {
			return nil, err
		}
	case <-time.After(DefaultTimeout):
		return nil, errors.New("timeout")
	}

	if err := c.clearDeadline(); err != nil {
		return nil, err
	}

	if cmd.response != nil {
		if err := DecodeResponse(c.res, cmd.response); err != nil {
			return nil, err
		}
	}

	return c.res, nil
}

// IsConnected returns whether the client is currently
// connected and processing incoming messages.
func (c *Client) IsConnected() bool {
	return len(c.disconnect) == 0
}

// Close closes the connection to the server.
func (c *Client) Close() error {
	if c.notify != nil {
		defer close(c.notify)
	}

	_, err := c.Exec("quit")
	err2 := c.conn.Close()

	if err != nil && err != ErrNotConnected {
		return err
	}

	return err2
}

// scanError returns the error from the scanner if non-nil,
// `io.ErrUnexpectedEOF` otherwise.
func (c *Client) scanErr() error {
	if err := c.scanner.Err(); err != nil {
		return err
	}
	return io.ErrUnexpectedEOF
}
