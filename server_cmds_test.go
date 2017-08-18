package ts3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCmdsServer(t *testing.T) {
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

	list := func(t *testing.T) {
		servers, err := c.Server.List()
		if !assert.NoError(t, err) {
			return
		}
		expected := []*Server{
			{
				ID:                 1,
				Port:               10677,
				Status:             "online",
				ClientsOnline:      1,
				QueryClientsOnline: 1,
				MaxClients:         35,
				Uptime:             12345025,
				Name:               "Server #1",
				AutoStart:          true,
				MachineID:          "1",
				UniqueIdentifier:   "uniq1",
			},
			{
				ID:                 2,
				Port:               10617,
				Status:             "online",
				ClientsOnline:      3,
				QueryClientsOnline: 2,
				MaxClients:         10,
				Uptime:             3165117,
				Name:               "Server #2",
				AutoStart:          true,
				MachineID:          "1",
				UniqueIdentifier:   "uniq2",
			},
		}
		assert.Equal(t, expected, servers)
	}

	idgetbyport := func(t *testing.T) {
		id, err := c.Server.IDGetByPort(1234)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, 1, id)
	}

	info := func(t *testing.T) {
		s, err := c.Server.Info()
		if !assert.NoError(t, err) {
			return
		}
		expected := &Server{
			Status:     "template",
			MaxClients: 32,
			Name:       "Test Server",
			AnitFloodPointsNeededCommandBlock:      150,
			AnitFloodPointsNeededIPBlock:           250,
			AnitFloodPointsTickReduce:              5,
			ComplainAutoBanCount:                   5,
			ComplainAutoBanTime:                    1200,
			ComplainRemoveTime:                     3600,
			DefaultChannelAdminGroup:               1,
			DefaultChannelGroup:                    4,
			DefaultServerGroup:                     5,
			MinClientsInChannelBeforeForcedSilence: 100,
			NeededIdentitySecurityLevel:            8,
			LogPermissions:                         true,
			PrioritySpeakerDimmModificator:         -18,
			MaxDownloadTotalBandwidth:              18446744073709551615,
			MaxUploadTotalBandwidth:                18446744073709551615,
			FileBase:                               "files",
			HostButtonToolTip:                      "Multiplay Game Servers",
			HostButtonURL:                          "http://www.multiplaygameservers.com",
			WelcomeMessage:                         "Welcome to TeamSpeak, check [URL]www.teamspeak.com[/URL] for latest infos.",
		}
		assert.Equal(t, expected, s)
	}

	create := func(t *testing.T) {
		s, err := c.Server.Create("my server")
		if !assert.NoError(t, err) {
			return
		}
		expected := &CreatedServer{
			ID:    2,
			Port:  9988,
			Token: "eKnFZQ9EK7G7MhtuQB6+N2B1PNZZ6OZL3ycDp2OW",
		}
		assert.Equal(t, expected, s)
	}

	edit := func(t *testing.T) {
		assert.NoError(t, c.Server.Edit(NewArg("virtualserver_maxclients", 10)))
	}

	del := func(t *testing.T) {
		assert.NoError(t, c.Server.Delete(1))
	}

	start := func(t *testing.T) {
		assert.NoError(t, c.Server.Start(1))
	}

	stop := func(t *testing.T) {
		assert.NoError(t, c.Server.Stop(1))
	}

	grouplist := func(t *testing.T) {
		groups, err := c.Server.GroupList()
		if !assert.NoError(t, err) {
			return
		}

		expected := []*Group{
			{
				ID:   1,
				Name: "Guest Server Query",
				Type: 2,
			},
			{
				ID:                2,
				Name:              "Admin Server Query",
				Type:              2,
				IconID:            500,
				Saved:             true,
				ModifyPower:       100,
				MemberAddPower:    100,
				MemberRemovePower: 100,
			},
		}
		assert.Equal(t, expected, groups)
	}

	privilegekeylist := func(t *testing.T) {
		keys, err := c.Server.PrivilegeKeyList()
		if !assert.NoError(t, err) {
			return
		}

		expected := []*PrivilegeKey{
			{
				Token:   "zTfamFVhiMEzhTl49KrOVYaMilHPDQEBQOJFh6qX",
				ID1:     17395,
				Created: 1499948005,
			},
		}
		assert.Equal(t, expected, keys)
	}

	privilegekeyadd := func(t *testing.T) {
		token, err := c.Server.PrivilegeKeyAdd(0, 17395, 0)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "zTfamFVhiMEzhTl49KrOVYaMilHPDQEBQOJFh6qX", token)
	}

	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{"list", list},
		{"idgetbyport", idgetbyport},
		{"info", info},
		{"create", create},
		{"edit", edit},
		{"del", del},
		{"start", start},
		{"stop", stop},
		{"grouplist", grouplist},
		{"privilegekeylist", privilegekeylist},
		{"privilegekeyadd", privilegekeyadd},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.f)
	}
}
