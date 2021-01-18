package gateway

const (
	// Gateway OP codes

	OP_DISPATCH = 0
	OP_HEARTBEAT = 1
	OP_IDENTIFY = 2
	OP_STATUS_UPDATE = 3

	// voice connection / disconnection
	OP_VOICE_UPDATE = 4
	OP_VOICE_PING = 5

	OP_RESUME = 6
	OP_RECONNECT = 7
	OP_REQ_GUILD_MEMBERS = 8
	OP_INVALID_SESSION = 9

	OP_HELLO = 10
	OP_HEARTBEAT_ACK = 11

	// request member / presence information
	OP_GUILD_SYNC = 12

	// request to sync up call dm / group dm
	OP_CALL_SYNC = 13

	// request for lazy guilds
	OP_LAZY_REQUEST = 14
	OP_UNKNOWN_1 = 23

)