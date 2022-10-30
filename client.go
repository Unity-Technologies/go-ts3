package ts3

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
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

	// DefaultConnectHeader send by server on connect.
	DefaultConnectHeader = "TS3"

	// startBufSize is the initial size of allocation for the parse buffer.
	startBufSize = 4096

	// keepAliveData is the keepalive data.
	keepAliveData = " \n"
)

var (
	respTrailerRe = regexp.MustCompile(`^error id=(\d+) msg=([^ ]+)(.*)`)

	// DefaultTimeout is the default read / write / dial timeout for Clients.
	DefaultTimeout = 10 * time.Second

	// DefaultKeepAlive is the default interval in which keepalive data is sent.
	DefaultKeepAlive = 200 * time.Second

	// DefaultNotifyBufSize is the default notification buffer size.
	DefaultNotifyBufSize = 5
)

// Client is a TeamSpeak 3 ServerQuery client.
type Client struct {
	conn          net.Conn
	timeout       time.Duration
	keepAlive     time.Duration
	scanner       *bufio.Scanner
	buf           []byte
	maxBufSize    int
	notifyBufSize int
	work          chan string
	err           chan error
	notify        chan Notification
	disconnect    chan struct{}
	res           []string
	connectHeader string

	Server *ServerMethods
}

// Timeout sets read / write / dial timeout for a TeamSpeak 3 Client.
func Timeout(timeout time.Duration) func(*Client) error {
	return func(c *Client) error {
		c.timeout = timeout
		return nil
	}
}

// KeepAlive sets the keepAlive interval.
func KeepAlive(keepAlive time.Duration) func(*Client) error {
	return func(c *Client) error {
		c.keepAlive = keepAlive
		return nil
	}
}

// NotificationBuffer sets the notification buffer size.
func NotificationBuffer(size int) func(*Client) error {
	return func(c *Client) error {
		c.notifyBufSize = size
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

// ConnectHeader sets the header expected on connect.
//
// Default is "TS3" which is sent by server query. For client query
// use "TS3 Client".
func ConnectHeader(connectHeader string) func(*Client) error {
	return func(c *Client) error {
		c.connectHeader = connectHeader
		return nil
	}
}

// NewClient returns a new TeamSpeak 3 client connected to addr.
func NewClient(addr string, options ...func(c *Client) error) (*Client, error) {
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%v:%v", addr, DefaultPort)
	}

	c := &Client{
		timeout:       DefaultTimeout,
		keepAlive:     DefaultKeepAlive,
		buf:           make([]byte, startBufSize),
		maxBufSize:    MaxParseTokenSize,
		notifyBufSize: DefaultNotifyBufSize,
		work:          make(chan string),
		err:           make(chan error),
		disconnect:    make(chan struct{}),
		connectHeader: DefaultConnectHeader,
	}
	for _, f := range options {
		if f == nil {
			return nil, ErrNilOption
		}
		if err := f(c); err != nil {
			return nil, err
		}
	}

	c.notify = make(chan Notification, c.notifyBufSize)

	// Wire up command groups
	c.Server = &ServerMethods{Client: c}

	var err error
	if c.conn, err = net.DialTimeout("tcp", addr, c.timeout); err != nil {
		return nil, fmt.Errorf("client: dial timeout: %w", err)
	}

	c.scanner = bufio.NewScanner(bufio.NewReader(c.conn))
	c.scanner.Buffer(c.buf, c.maxBufSize)
	c.scanner.Split(ScanLines)

	if err := c.conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, fmt.Errorf("client: set deadline: %w", err)
	}

	// Read the connection header
	if !c.scanner.Scan() {
		return nil, c.scanErr()
	}

	if l := c.scanner.Text(); l != c.connectHeader {
		return nil, fmt.Errorf("client: invalid connection header %q", l)
	}

	// Slurp the banner
	if !c.scanner.Scan() {
		return nil, c.scanErr()
	}

	if err := c.conn.SetReadDeadline(time.Time{}); err != nil {
		return nil, fmt.Errorf("client: set read deadline: %w", err)
	}

	// Start handlers
	go c.messageHandler()
	go c.workHandler()

	return c, nil
}

// messageHandler scans incoming lines and handles them accordingly.
func (c *Client) messageHandler() {
	for {
		if c.scanner.Scan() {
			line := c.scanner.Text()
			//nolint: gocritic
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
		} else {
			err := c.scanErr()
			c.err <- err
			if errors.Is(err, io.ErrUnexpectedEOF) {
				close(c.disconnect)
				return
			}
		}
	}
}

// workHandler handles commands and keepAlive messages.
func (c *Client) workHandler() {
	for {
		select {
		case w := <-c.work:
			c.process(w)
		case <-time.After(c.keepAlive):
			c.process(keepAliveData)
		case <-c.disconnect:
			return
		}
	}
}

func (c *Client) process(data string) {
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		c.err <- err
	}
	if _, err := c.conn.Write([]byte(data)); err != nil {
		c.err <- err
	}
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

	c.work <- cmd.String()

	select {
	case err := <-c.err:
		if err != nil {
			return nil, err
		}
	case <-time.After(c.timeout):
		return nil, ErrTimeout
	}

	res := c.res
	c.res = nil

	if cmd.response != nil {
		if err := DecodeResponse(res, cmd.response); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// IsConnected returns whether the client is connected.
func (c *Client) IsConnected() bool {
	select {
	case <-c.disconnect:
		return false
	default:
		return true
	}
}

// Close closes the connection to the server.
func (c *Client) Close() error {
	defer close(c.notify)

	_, err := c.Exec("quit")
	err2 := c.conn.Close()

	if err != nil {
		return err
	} else if err2 != nil {
		return fmt.Errorf("client: close: %w", err2)
	}

	return nil
}

// scanError returns the error from the scanner if non-nil,
// `io.ErrUnexpectedEOF` otherwise.
func (c *Client) scanErr() error {
	if err := c.scanner.Err(); err != nil {
		return fmt.Errorf("client: scan: %w", err)
	}
	return io.ErrUnexpectedEOF
}
