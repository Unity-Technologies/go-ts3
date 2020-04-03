package ts3

// Login authenticates with the server.
func (c *Client) Login(user, passwd string) error {
	_, err := c.ExecCmd(NewCmd("login").WithArgs(
		NewArg("client_login_name", user),
		NewArg("client_login_password", passwd)),
	)
	return err
}

// Logout deselect virtual server and log out.
func (c *Client) Logout() error {
	_, err := c.Exec("logout")
	return err
}

// Version represents version information.
type Version struct {
	Version  string
	Platform string
	Build    int
}

// Version returns version information.
func (c *Client) Version() (*Version, error) {
	v := &Version{}
	if _, err := c.ExecCmd(NewCmd("version").WithResponse(v)); err != nil {
		return nil, err
	}

	return v, nil
}

// Use selects a virtual server by id.
func (c *Client) Use(id int) error {
	_, err := c.ExecCmd(NewCmd("use").WithArgs(NewArg("sid", id)))
	return err
}

// UsePort selects a virtual server by port.
func (c *Client) UsePort(port int) error {
	_, err := c.ExecCmd(NewCmd("use").WithArgs(NewArg("port", port)))
	return err
}

// ConnectionInfo represents an answer of the whoami command.
type ConnectionInfo struct {
	ServerStatus           string `ms:"virtualserver_status"`
	ServerID               int    `ms:"virtualserver_id"`
	ServerUniqueIdentifier string `ms:"virtualserver_unique_identifier"`
	ServerPort             uint16 `ms:"virtualserver_port"`
	ClientID               int    `ms:"client_id"`
	ClientChannelID        int    `ms:"client_channel_id"`
	ClientName             string `ms:"client_nickname"`
	ClientDatabaseID       int    `ms:"client_database_id"`
	ClientLoginName        string `ms:"client_login_name"`
	ClientUniqueIdentifier string `ms:"client_unique_identifier"`
	ClientOriginServerID   int    `ms:"client_origin_server_id"`
}

// Whoami returns information about the current connection including the currently selected virtual server.
func (c *Client) Whoami() (*ConnectionInfo, error) {
	i := &ConnectionInfo{}
	if _, err := c.ExecCmd(NewCmd("whoami").WithResponse(&i)); err != nil {
		return nil, err
	}

	return i, nil
}

// Properties that can be changed with ClientUpdate
const (
	ClientNickname           = "client_nickname"
	ClientIsTalker           = "client_is_talker"
	ClientDescription        = "client_description"
	ClientIsChannelCommander = "client_is_channel_commander"
	ClientIconID             = "client_icon_id"
)

// ClientUpdate changes properties of the client to a given value.
func (c *Client) ClientUpdate(properties ...CmdArg) error {
	_, err := c.ExecCmd(NewCmd("clientupdate").WithArgs(properties...))
	return err
}

// SetNick sets the clients nickname.
func (c *Client) SetNick(nick string) error {
	return c.ClientUpdate(NewArg(ClientNickname, nick))
}

// SetTalker sets whether the client is able to talk.
func (c *Client) SetTalker(val bool) error {
	return c.ClientUpdate(NewArg(ClientIsTalker, val))
}

// SetDescription sets the clients description.
func (c *Client) SetDescription(description string) error {
	return c.ClientUpdate(NewArg(ClientDescription, description))
}

// SetChannelCommander sets whether the client is a channel commander.
func (c *Client) SetChannelCommander(val bool) error {
	return c.ClientUpdate(NewArg(ClientIsChannelCommander, val))
}

// SetIcon sets the clients icon based on the CRC32 checksum.
func (c *Client) SetIcon(id int) error {
	return c.ClientUpdate(NewArg(ClientIconID, id))
}
