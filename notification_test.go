package ts3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeNotification(t *testing.T) {
	r := decodeNotification(`notifytextmessage targetmode=3 msg=lorem\sipsum invokerid=42 invokername=foobar invokeruid=something= flag`)
	expected := Notification{
		Type: "textmessage",
		Data: map[string]string{
			"targetmode":  "3",
			"msg":         "lorem ipsum",
			"invokerid":   "42",
			"invokername": "foobar",
			"invokeruid":  "something=",
			"flag":        "",
		},
	}
	assert.Equal(t, expected, r)

}
