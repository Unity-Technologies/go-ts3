package ts3

import (
	"errors"
	"strings"
)

// Notify event types
// Notifications are very badly documented by TeamSpeak;
// A complete but unofficial documentation in German can be found here:
// http://yat.qa/ressourcen/server-query-notify/
const (
	// ServerEvents registers the following events:
	// `cliententerview`, `clientleftview`, `serveredited`
	ServerEvents = "server"

	// ServerEvents registers the following events:
	// `cliententerview`, `clientleftview`, `channeldescriptionchanged`, `channelpasswordchanged`
	// `channelmoved`, `channeledited`, `channelcreated`, `channeldeleted`, `clientmoved`
	ChannelEvents = "channel"

	// TextServerEvents registers the `textmessage` event with `targetmode = 3`
	TextServerEvents = "textserver"

	// TextChannelEvents registers the `textmessage` event with `targetmode = 2`
	//
	// Notifications are only received for messages that are sent in the channel that the client is in.
	TextChannelEvents = "textchannel"

	// TextPrivateEvents registers the `textmessage` event with `targetmode = 1`
	TextPrivateEvents = "textprivate"

	// TokenUsedEvents registers the `tokenused` event
	TokenUsedEvents = "tokenused"
)

// Notification contains the information of a notify event
type Notification struct {
	Type string
	Data map[string]string
}

// SetNotifyHandler registers a func that handles received notifications
func (c *Client) SetNotifyHandler(event string, handler func(Notification)) {
	c.notifyHandler = handler
}

// NotifyRegister unregisters a given event
//
// The id can only be set for channel events
// It's not possible to subscribe to multiple channel id's;
// To receive events for all channels the id can be left out or set to 0.
func (c *Client) NotifyRegister(event string, id ...uint) error {
	if len(id) == 0 {
		if event == ChannelEvents {
			// if no channel ID is set use 0 (all channels) as default.
			return c.NotifyRegister(event, 0)
		}

		_, err := c.ExecCmd(NewCmd("servernotifyregister").WithArgs(
			NewArg("event", event),
		))
		return err
	} else if len(id) == 1 {
		_, err := c.ExecCmd(NewCmd("servernotifyregister").WithArgs(
			NewArg("event", event),
			NewArg("id", id),
		))
		return err
	} else {
		return errors.New("invalid Argument, only one ID can be set")
	}
}

// NotifyUnregister unregisters a given event
func (c *Client) NotifyUnregister(event string) error {
	_, err := c.ExecCmd(NewCmd("servernotifyunregister").WithArgs(NewArg("event", event)))
	return err
}

func (c *Client) notifyDispatcher() {
	for {
		notification := <-c.notify
		if c.notifyHandler != nil {
			c.notifyHandler(decodeNotification(notification))
		}
	}
}

func decodeNotification(str string) Notification {
	params := strings.Split(str, " ")
	data := map[string]string{}

	for _, param := range params[1:] {
		kvPair := strings.SplitN(param, "=", 2)
		if len(kvPair) > 1 {
			data[kvPair[0]] = Decode(kvPair[1])
		} else {
			data[kvPair[0]] = ""
		}
	}

	return Notification{
		Type: strings.Replace(params[0], "notify", "", 1),
		Data: data,
	}
}
