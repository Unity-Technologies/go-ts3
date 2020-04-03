package ts3

import (
	"time"
)

const (
	// ExtendedServerList can be passed to List to get extended server information.
	ExtendedServerList = "-extended"
)

// ServerMethods groups server methods.
type ServerMethods struct {
	*Client
}

// Instance represents basic information for a TeamSpeak 3 instance.
// This is the result of an instanceinfo call.
type Instance struct {
	DatabaseVersion                int    `ms:"serverinstance_database_version"`
	FileTransferPort               uint16 `ms:"serverinstance_filetransfer_port"`
	MaxTotalDownloadBandwidth      uint64 `ms:"serverinstance_max_download_total_bandwidth"`
	MaxTotalUploadBandwidth        uint64 `ms:"serverinstance_max_upload_total_bandwidth"`
	GuestServerQueryGroup          int    `ms:"serverinstance_guest_serverquery_group"`
	ServerQueryFloodCommands       int    `ms:"serverinstance_serverquery_flood_commands"`
	ServerQueryFloodTime           int    `ms:"serverinstance_serverquery_flood_time"`
	ServerQueryBanTime             int    `ms:"serverinstance_serverquery_ban_time"`
	TemplateServerAdminGroup       int    `ms:"serverinstance_template_serveradmin_group"`
	TemplateServerDefaultGroup     int    `ms:"serverinstance_template_serverdefault_group"`
	TemplateChannelAdminGroup      int    `ms:"serverinstance_template_channeladmin_group"`
	TemplateChannelDefaultGroup    int    `ms:"serverinstance_template_channeldefault_group"`
	PermissionsVersion             int    `ms:"serverinstance_permissions_version"`
	PendingConnectionsPerIP        int    `ms:"serverinstance_pending_connections_per_ip"`
	ServerQueryMaxConnectionsPerIP int    `ms:"serverinstance_serverquery_max_connections_per_ip"`
}

// ServerConnectionInfo represents the connection info for a TeamSpeak 3 instance.
// This is the result of an hostinfo call.
type ServerConnectionInfo struct {
	InstanceUptime                    time.Duration `ms:"instance_uptime"`
	HostTimestamp                     int64         `ms:"host_timestamp_utc"`
	VirtualServersRunning             uint          `ms:"virtualservers_running_total"`
	VirtualServersTotalMaxClients     uint16        `ms:"virtualservers_total_maxclients"`
	VirtualServersTotalClientsOnline  uint16        `ms:"virtualservers_total_clients_online"`
	VirtualServersTotalChannelsOnline uint16        `ms:"virtualservers_total_channels_online"`
	FileTransferBandwidthSent         uint64        `ms:"connection_filetransfer_bandwidth_sent"`
	FileTransferBandwidthReceived     uint64        `ms:"connection_filetransfer_bandwidth_received"`
	FileTransferBytesSentTotal        uint64        `ms:"connection_filetransfer_bytes_sent_total"`
	FileTransferBytesReceivedTotal    uint64        `ms:"connection_filetransfer_bytes_received_total"`
	PacketsSentTotal                  uint64        `ms:"connection_packets_sent_total"`
	PacketsReceivedTotal              uint64        `ms:"connection_packets_received_total"`
	BytesSentTotal                    uint64        `ms:"connection_bytes_sent_total"`
	BytesReceivedTotal                uint64        `ms:"connection_bytes_received_total"`
	BandwidthSentLastSecond           uint64        `ms:"connection_bandwidth_sent_last_second_total"`
	BandwidthReceivedLastSecond       uint64        `ms:"connection_bandwidth_received_last_second_total"`
	BandwidthSentLastMinute           uint64        `ms:"connection_bandwidth_sent_last_minute_total"`
	BandwidthReceivedLastMinute       uint64        `ms:"connection_bandwidth_received_last_minute_total"`
}

// Server represents a TeamSpeak 3 virtual server.
// This is the result of an serverinfo call.
type Server struct {
	UniqueIdentifier                       string        `ms:"virtualserver_unique_identifier"`
	Name                                   string        `ms:"virtualserver_name"`
	WelcomeMessage                         string        `ms:"virtualserver_welcomemessage"`
	Platform                               string        `ms:"virtualserver_platform"`
	Version                                string        `ms:"virtualserver_version"`
	Password                               string        `ms:"virtualserver_password"`
	HostMessage                            string        `ms:"virtualserver_hostmessage"`
	FileBase                               string        `ms:"virtualserver_filebase"`
	HostBannerURL                          string        `ms:"virtualserver_hostbanner_url"`
	HostBannerGFXURL                       string        `ms:"virtualserver_hostbanner_gfx_url"`
	HostButtonToolTip                      string        `ms:"virtualserver_hostbutton_tooltip"`
	HostButtonURL                          string        `ms:"virtualserver_hostbutton_url"`
	HostButtonGFXURL                       string        `ms:"virtualserver_hostbutton_gfx_url"`
	MachineID                              string        `ms:"virtualserver_machine_id"`
	NamePhonetic                           string        `ms:"virtualserver_name_phonetic"`
	IP                                     string        `ms:"virtualserver_ip"`
	ServerNickname                         string        `ms:"virtualserver_nickname"`
	Status                                 string        `ms:"virtualserver_status"`
	Created                                int64         `ms:"virtualserver_created"`
	MaxDownloadTotalBandwidth              uint64        `ms:"virtualserver_max_download_total_bandwidth"`
	MaxUploadTotalBandwidth                uint64        `ms:"virtualserver_max_upload_total_bandwidth"`
	ClientConnections                      uint64        `ms:"virtualserver_client_connections"`
	QueryClientConnections                 uint64        `ms:"virtualserver_query_client_connections"`
	DownloadQuota                          uint64        `ms:"virtualserver_download_quota"`
	UploadQuota                            uint64        `ms:"virtualserver_upload_quota"`
	FileTransferBandwidthSent              uint64        `ms:"connection_filetransfer_bandwidth_sent"`
	FileTransferBandwidthReceived          uint64        `ms:"connection_filetransfer_bandwidth_received"`
	FileTransferBytesSentTotal             uint64        `ms:"connection_filetransfer_bytes_sent_total"`
	FileTransferBytesReceivedTotal         uint64        `ms:"connection_filetransfer_bytes_received_total"`
	PacketsSentSpeech                      uint64        `ms:"connection_packets_sent_speech"`
	BytesSentSpeech                        uint64        `ms:"connection_bytes_sent_speech"`
	PacketsReceivedSpeech                  uint64        `ms:"connection_packets_received_speech"`
	BytesReceivedSpeech                    uint64        `ms:"connection_bytes_received_speech"`
	PacketsSentKeepalive                   uint64        `ms:"connection_packets_sent_keepalive"`
	BytesSentKeepalive                     uint64        `ms:"connection_bytes_sent_keepalive"`
	PacketsReceivedKeepalive               uint64        `ms:"connection_packets_received_keepalive"`
	BytesReceivedKeepalive                 uint64        `ms:"connection_bytes_received_keepalive"`
	PacketsSentControl                     uint64        `ms:"connection_packets_sent_control"`
	BytesSentControl                       uint64        `ms:"connection_bytes_sent_control"`
	PacketsReceivedControl                 uint64        `ms:"connection_packets_received_control"`
	BytesReceivedControl                   uint64        `ms:"connection_bytes_received_control"`
	PacketsSentTotal                       uint64        `ms:"connection_packets_sent_total"`
	BytesSentTotal                         uint64        `ms:"connection_bytes_sent_total"`
	PacketsReceivedTotal                   uint64        `ms:"connection_packets_received_total"`
	BytesReceivedTotal                     uint64        `ms:"connection_bytes_received_total"`
	BandwidthSentLastSecondTotal           uint64        `ms:"virtualserver_bandwidth_sent_last_second_total"`
	BandwidthSentLastMinuteTotal           uint64        `ms:"virtualserver_bandwidth_sent_last_minute_total"`
	BandwidthReceivedLastSecondTotal       uint64        `ms:"virtualserver_bandwidth_received_last_second_total"`
	BandwidthReceivedLastMinuteTotal       uint64        `ms:"virtualserver_bandwidth_received_last_minute_total"`
	MonthBytesDownloaded                   uint64        `ms:"virtualserver_month_bytes_downloaded"`
	MonthBytesUploaded                     uint64        `ms:"virtualserver_month_bytes_uploaded"`
	TotalBytesDownloaded                   uint64        `ms:"virtualserver_total_bytes_downloaded"`
	TotalBytesUploaded                     uint64        `ms:"virtualserver_total_bytes_uploaded"`
	TotalPacketLossSpeech                  float64       `ms:"virtualserver_total_packetloss_speech"`
	TotalPacketLossKeepalive               float64       `ms:"virtualserver_total_packetloss_keepalive"`
	TotalPacketLossControl                 float64       `ms:"virtualserver_total_packetloss_control"`
	TotalPacketLossTotal                   float64       `ms:"virtualserver_total_packetloss_total"`
	TotalPing                              float32       `ms:"virtualserver_total_ping"`
	PrioritySpeakerDimmModificator         float32       `ms:"virtualserver_priority_speaker_dimm_modificator"`
	MaxClients                             uint16        `ms:"virtualserver_maxclients"`
	ClientsOnline                          uint16        `ms:"virtualserver_clientsonline"`
	ChannelsOnline                         uint16        `ms:"virtualserver_channelsonline"`
	Uptime                                 time.Duration `ms:"virtualserver_uptime"`
	CodecEncryptionMode                    uint          `ms:"virtualserver_codec_encryption_mode"`
	HostMessageMode                        uint          `ms:"virtualserver_hostmessage_mode"`
	DefaultServerGroup                     uint          `ms:"virtualserver_default_server_group"`
	DefaultChannelGroup                    uint          `ms:"virtualserver_default_channel_group"`
	FlagPassword                           uint          `ms:"virtualserver_flag_password"`
	DefaultChannelAdminGroup               uint          `ms:"virtualserver_default_channel_admin_group"`
	HostBannerGFXInterval                  uint          `ms:"virtualserver_hostbanner_gfx_interval"`
	ComplainAutoBanCount                   uint          `ms:"virtualserver_complain_autoban_count"`
	ComplainAutoBanTime                    uint          `ms:"virtualserver_complain_autoban_time"`
	ComplainRemoveTime                     uint          `ms:"virtualserver_complain_remove_time"`
	MinClientsInChannelBeforeForcedSilence uint          `ms:"virtualserver_min_clients_in_channel_before_forced_silence"`
	ID                                     int           `ms:"virtualserver_id"`
	AntiFloodPointsTickReduce              int           `ms:"virtualserver_antiflood_points_tick_reduce"`
	AntiFloodPointsNeededCommandBlock      int           `ms:"virtualserver_antiflood_points_needed_command_block"`
	AntiFloodPointsNeededIPBlock           int           `ms:"virtualserver_antiflood_points_needed_ip_block"`
	QueryClientsOnline                     uint16        `ms:"virtualserver_queryclientsonline"`
	Port                                   uint16        `ms:"virtualserver_port"`
	AutoStart                              int           `ms:"virtualserver_autostart"`
	NeededIdentitySecurityLevel            int           `ms:"virtualserver_needed_identity_security_level"`
	LogClient                              int           `ms:"virtualserver_log_client"`
	LogQuery                               int           `ms:"virtualserver_log_client"`
	LogChannel                             int           `ms:"virtualserver_log_channel"`
	LogPermissions                         int           `ms:"virtualserver_log_permissions"`
	LogServer                              int           `ms:"virtualserver_log_server"`
	LogFileTransfer                        int           `ms:"virtualserver_log_filetransfer"`
	MinClientVersion                       int           `ms:"virtualserver_min_client_version"`
	IconID                                 int           `ms:"virtualserver_icon_id"`
	ReservedSlots                          int           `ms:"virtualserver_reserved_slots"`
	WebListEnabled                         int           `ms:"virtualserver_web_list_enabled"`
	AskForPrivilegeKey                     int           `ms:"virtualserver_ask_for_privilegekey"`
	HostBannerMode                         int           `ms:"virtualserver_hostbanner_mode"`
	ChannelTempDeleteDelayDefault          int           `ms:"virtualserver_channel_temp_delete_delay_default"`
	MinAndroidVersion                      int           `ms:"virtualserver_min_android_version"`
	MiniOSVersion                          int           `ms:"virtualserver_min_ios_version"`
	AntiFloodPointsNeededPluginBlock       int           `ms:"virtualserver_antiflood_points_needed_plugin_block"`
}

// List lists virtual servers.
// In addition to the options supported by the Teamspeak 3 query protocol it also supports the ExtendedServerList option.
// If ExtendedServerList is specified in options then each server returned contain extended server information as returned by Info.
func (s *ServerMethods) List(options ...string) (servers []*Server, err error) {
	var extended bool
	for i, o := range options {
		if o == ExtendedServerList {
			options = append(options[:i], options[i+1:]...)
			extended = true
		}
	}

	if _, err = s.ExecCmd(NewCmd("serverlist").WithOptions(options...).WithResponse(&servers)); err != nil {
		return nil, err
	}

	if extended {
		var info *ConnectionInfo
		if info, err = s.Whoami(); err != nil {
			return nil, err
		}

		var lastID int
		defer func() {
			if lastID != info.ServerID {
				// Restore the previously selected server
				if err2 := s.Use(info.ServerID); err2 != nil && err != nil {
					err = err2
				}
			}
		}()

		for _, server := range servers {
			if err = s.Use(server.ID); err != nil {
				return nil, err
			}
			lastID = server.ID

			if _, err = s.ExecCmd(NewCmd("serverinfo").WithResponse(server)); err != nil {
				return nil, err
			}
		}
	}

	return servers, nil
}

// IDGetByPort returns the database id of the virtual server running on UDP port.
func (s *ServerMethods) IDGetByPort(port uint16) (int, error) {
	r := struct {
		ID int `ms:"server_id"`
	}{}
	_, err := s.ExecCmd(NewCmd("serveridgetbyport").WithArgs(NewArg("virtualserver_port", port)).WithResponse(&r))
	return r.ID, err
}

// Info returns detailed configuration information about the selected server.
func (s *ServerMethods) Info() (*Server, error) {
	r := &Server{}
	if _, err := s.ExecCmd(NewCmd("serverinfo").WithResponse(&r)); err != nil {
		return nil, err
	}

	return r, nil
}

// InstanceInfo returns detailed information about the selected instance.
func (s *ServerMethods) InstanceInfo() (*Instance, error) {
	r := &Instance{}
	if _, err := s.ExecCmd(NewCmd("instanceinfo").WithResponse(&r)); err != nil {
		return nil, err
	}

	return r, nil
}

// ServerConnectionInfo returns detailed bandwidth and transfer information about the selected instance.
func (s *ServerMethods) ServerConnectionInfo() (*ServerConnectionInfo, error) {
	r := &ServerConnectionInfo{}
	if _, err := s.ExecCmd(NewCmd("serverrequestconnectioninfo").WithResponse(&r)); err != nil {
		return nil, err
	}

	return r, nil
}

// Edit changes the selected virtual servers configuration using the given args.
func (s *ServerMethods) Edit(args ...CmdArg) error {
	_, err := s.ExecCmd(NewCmd("serveredit").WithArgs(args...))
	return err
}

// Delete deletes the virtual server specified by id.
// Only virtual server in a stopped state can be deleted.
func (s *ServerMethods) Delete(id int) error {
	_, err := s.ExecCmd(NewCmd("serverdelete").WithArgs(NewArg("sid", id)))
	return err
}

// CreatedServer is the details returned by a server create.
type CreatedServer struct {
	ID    int    `ms:"sid"`
	Port  uint16 `ms:"virtualserver_port"`
	Token string
}

// Create creates a new virtual server using the given properties and returns
// its ID, port and initial administrator privilege key.
// If virtualserver_port arg is not specified, the server will use the first unused
// UDP port.
func (s *ServerMethods) Create(name string, args ...CmdArg) (*CreatedServer, error) {
	r := &CreatedServer{}
	args = append(args, NewArg("virtualserver_name", name))
	if _, err := s.ExecCmd(NewCmd("servercreate").WithArgs(args...).WithResponse(r)); err != nil {
		return nil, err
	}

	return r, nil
}

// Start starts the virtual server specified by id.
func (s *ServerMethods) Start(id int) error {
	_, err := s.ExecCmd(NewCmd("serverstart").WithArgs(NewArg("sid", id)))
	return err
}

// Stop stops the virtual server specified by id.
func (s *ServerMethods) Stop(id int) error {
	_, err := s.ExecCmd(NewCmd("serverstop").WithArgs(NewArg("sid", id)))
	return err
}

// Group represents a virtual server group.
type Group struct {
	ID                int `ms:"sgid"`
	Name              string
	Type              int
	IconID            int
	Saved             bool `ms:"savedb"`
	SortID            int
	NameMode          int
	ModifyPower       int `ms:"n_modifyp"`
	MemberAddPower    int `ms:"n_member_addp"`
	MemberRemovePower int `ms:"n_member_addp"`
}

// GroupList returns a list of available groups for the selected server.
func (s *ServerMethods) GroupList() ([]*Group, error) {
	var groups []*Group
	if _, err := s.ExecCmd(NewCmd("servergrouplist").WithResponse(&groups)); err != nil {
		return nil, err
	}

	return groups, nil
}

// Channel represents a TeamSpeak 3 channel in a virtual server.
type Channel struct {
	ID                   int    `ms:"cid"`
	ParentID             int    `ms:"pid"`
	ChannelOrder         int    `ms:"channel_order"`
	ChannelName          string `ms:"channel_name"`
	TotalClients         int    `ms:"total_clients"`
	NeededSubscribePower int    `ms:"channel_needed_subscribe_power"`
}

// ChannelList returns a list of channels for the selected server.
func (s *ServerMethods) ChannelList() ([]*Channel, error) {
	var channels []*Channel
	if _, err := s.ExecCmd(NewCmd("channellist").WithResponse(&channels)); err != nil {
		return nil, err
	}

	return channels, nil
}

// PrivilegeKey represents a server privilege key.
type PrivilegeKey struct {
	Token       string
	Type        int    `ms:"token_type"`
	ID1         int    `ms:"token_id1"`
	ID2         int    `ms:"token_id2"`
	Created     int    `ms:"token_created"`
	Description string `ms:"token_description"`
}

// PrivilegeKeyList returns a list of available privilege keys for the selected server,
// including their type and group IDs.
func (s *ServerMethods) PrivilegeKeyList() ([]*PrivilegeKey, error) {
	var keys []*PrivilegeKey
	if _, err := s.ExecCmd(NewCmd("privilegekeylist").WithResponse(&keys)); err != nil {
		return nil, err
	}

	return keys, nil
}

// PrivilegeKeyAdd creates a new privilege token to the selected server and returns it.
// If tokentype is set to 0, the ID specified with id1 will be a server group ID.
// Otherwise, id1 is used as a channel group ID and you need to provide a valid channel ID using id2.
func (s *ServerMethods) PrivilegeKeyAdd(ttype, id1, id2 int, options ...CmdArg) (string, error) {
	t := struct {
		Token string
	}{}
	options = append(options, NewArg("tokentype", ttype), NewArg("tokenid1", id1), NewArg("tokenid2", id2))
	_, err := s.ExecCmd(NewCmd("privilegekeylist").WithArgs(options...).WithResponse(&t))
	return t.Token, err
}

// OnlineClient represents a client online on a virtual server.
type OnlineClient struct {
	ID          int    `ms:"clid"`
	ChannelID   int    `ms:"cid"`
	DatabaseID  int    `ms:"client_database_id"`
	Nickname    string `ms:"client_nickname"`
	Type        int    `ms:"client_type"`
	Away        bool   `ms:"client_away"`
	AwayMessage string `ms:"client_away_message"`
}

// ClientList returns a list of online clients.
func (s *ServerMethods) ClientList() ([]*OnlineClient, error) {
	var clients []*OnlineClient
	if _, err := s.ExecCmd(NewCmd("clientlist").WithResponse(&clients)); err != nil {
		return nil, err
	}
	return clients, nil
}

// DBClient represents a client identity on a virtual server.
type DBClient struct {
	ID               int       `ms:"cldbid"`
	UniqueIdentifier string    `ms:"client_unique_identifier"`
	Nickname         string    `ms:"client_nickname"`
	Created          time.Time `ms:"client_created"`
	LastConnected    time.Time `ms:"client_lastconnected"`
	Connections      int       `ms:"client_totalconnections"`
}

// ClientDBList returns a list of client identities known by the server.
func (s *ServerMethods) ClientDBList() ([]*DBClient, error) {
	var dbclients []*DBClient
	if _, err := s.ExecCmd(NewCmd("clientdblist").WithResponse(&dbclients)); err != nil {
		return nil, err
	}
	return dbclients, nil
}
