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
