package ts3

import (
	"strings"
)

// NotifyEvent is an event type.
type NotifyEvent string

const (
	// ServerEvents registers the following events:
	// `cliententerview`, `clientleftview`, `serveredited`.
	ServerEvents NotifyEvent = "server"

	// ChannelEvents registers the following events:
	// `cliententerview`, `clientleftview`, `channeldescriptionchanged`, `channelpasswordchanged`
	// `channelmoved`, `channeledited`, `channelcreated`, `channeldeleted`, `clientmoved`.
	ChannelEvents NotifyEvent = "channel"

	// TextServerEvents registers the `textmessage` event with `targetmode = 3`.
	TextServerEvents NotifyEvent = "textserver"

	// TextChannelEvents registers the `textmessage` event with `targetmode = 2`.
	//
	// Notifications are only received for messages that are sent in the channel that the client is in.
	TextChannelEvents NotifyEvent = "textchannel"

	// TextPrivateEvents registers the `textmessage` event with `targetmode = 1`.
	TextPrivateEvents NotifyEvent = "textprivate"

	// TokenUsedEvents registers the `tokenused` event.
	TokenUsedEvents NotifyEvent = "tokenused"
)

// Notification contains the information of a notify event.
type Notification struct {
	Type string
	Data map[string]string
}

// Notifications returns a read-only channel that outputs received notifications.
//
// If you subscribe to server and channel events you will receive duplicate
// `cliententerview` and `clientleftview` notifications.
// Sending a private message from the client results in a `textmessage`
// Notification even if the client didn't subscribe to any events.
//
// Notifications are not documented by TeamSpeak;
// A complete but unofficial documentation in German can be found here:
// http://yat.qa/ressourcen/server-query-notify/
func (c *Client) Notifications() <-chan Notification {
	return c.notify
}

// Register registers for server event notifications.
//
// Subscriptions will be reset on `logout`, `login` or `use`.
func (c *Client) Register(event NotifyEvent) error {
	if event == ChannelEvents {
		return c.RegisterChannel(0)
	}

	_, err := c.ExecCmd(NewCmd("servernotifyregister").WithArgs(
		NewArg("event", event),
	))
	return err
}

// RegisterChannel registers for channel event notifications.
//
// It's not possible to subscribe to multiple channels.
// To receive events for all channels the id can be set to 0.
func (c *Client) RegisterChannel(id uint) error {
	_, err := c.ExecCmd(NewCmd("servernotifyregister").WithArgs(
		NewArg("event", ChannelEvents),
		NewArg("id", id),
	))
	return err
}

// Unregister unregisters all events previously registered.
func (c *Client) Unregister() error {
	_, err := c.Exec("servernotifyunregister")
	return err
}

func decodeNotification(str string) (Notification, error) {
	parts := strings.SplitN(str, " ", 2)
	n := Notification{
		Type: strings.TrimPrefix(parts[0], "notify"),
	}

	err := DecodeResponse([]string{parts[1]}, &n.Data)

	return n, err
}
