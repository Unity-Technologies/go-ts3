package ts3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCmdsBasic(t *testing.T) {
	s := newServer(t)
	if s == nil {
		return
	}
	defer func() {
		assert.NoError(t, s.Close())
	}()

	c, err := NewClient(s.Addr, Timeout(time.Second*2))
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, c.Close())
	}()

	auth := func(t *testing.T) {
		if err = c.Login("user", "pass"); !assert.NoError(t, err) {
			return
		}

		if err = c.Logout(); !assert.NoError(t, err) {
			return
		}
	}

	version := func(t *testing.T) {
		v, err := c.Version()
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "3.0.12.2", v.Version)
		assert.Equal(t, 1455547898, v.Build)
		assert.Equal(t, "FreeBSD", v.Platform)
	}

	useID := func(t *testing.T) {
		assert.NoError(t, c.Use(1))
	}

	usePort := func(t *testing.T) {
		assert.NoError(t, c.UsePort(1024))
	}

	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{"auth", auth},
		{"version", version},
		{"useid", useID},
		{"userport", usePort},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.f)
	}
}
