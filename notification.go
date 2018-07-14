package ts3

import (
	"strings"
)

type notifyEvent string

// Notify event categorys
const (
	// ServerEvents registers the following events:
	// `cliententerview`, `clientleftview`, `serveredited`.
	ServerEvents notifyEvent = "server"

	// ServerEvents registers the following events:
	// `cliententerview`, `clientleftview`, `channeldescriptionchanged`, `channelpasswordchanged`
	// `channelmoved`, `channeledited`, `channelcreated`, `channeldeleted`, `clientmoved`.
	ChannelEvents notifyEvent = "channel"

	// TextServerEvents registers the `textmessage` event with `targetmode = 3`.
	TextServerEvents notifyEvent = "textserver"

	// TextChannelEvents registers the `textmessage` event with `targetmode = 2`.
	//
	// Notifications are only received for messages that are sent in the channel that the client is in.
	TextChannelEvents notifyEvent = "textchannel"

	// TextPrivateEvents registers the `textmessage` event with `targetmode = 1`.
	TextPrivateEvents notifyEvent = "textprivate"

	// TokenUsedEvents registers the `tokenused` event.
	TokenUsedEvents notifyEvent = "tokenused"
)

// notificationBuffer
const notificationBuffer = 10

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
	if c.notify == nil {
		c.notify = make(chan Notification, notificationBuffer)
	}
	return c.notify
}

// NotifyRegister subscribes for a specified category of events
// on a virtual server to receive notifications.
//
// Subscriptions will be reset on `logout`, `login` or `use`.
func (c *Client) NotifyRegister(event notifyEvent) error {
	if event == ChannelEvents {
		return c.NotifyRegisterChannel(0)
	}

	_, err := c.ExecCmd(NewCmd("servernotifyregister").WithArgs(
		NewArg("event", event),
	))
	return err
}

// NotifyRegisterChannel registers for events of a certain channel.
//
// It's not possible to subscribe to multiple channels.
// To receive events for all channels the id can be set to 0.
func (c *Client) NotifyRegisterChannel(id uint) error {
	_, err := c.ExecCmd(NewCmd("servernotifyregister").WithArgs(
		NewArg("event", ChannelEvents),
		NewArg("id", id),
	))
	return err
}

// NotifyUnregister unregisters all events previously registered
func (c *Client) NotifyUnregister() error {
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
