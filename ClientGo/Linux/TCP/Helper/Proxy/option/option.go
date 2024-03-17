package option

const (
	TCP_BUFFER_SIZE = 0x8000

	// UDP protocol's max capacity
	UDP_PACKET_MAX_SIZE = 0xFFFF - 28

	UDP_PACKET_CHANNEL_SIZE = 0x800

	CONNECTING_RETRY_DURATION = 1500

	SMUX_KEEPALIVE_INTERVAL = 20
	SMUX_KEEPALIVE_TIMEOUT  = 60
	SMUX_FRAMESIZE          = 0x8000
	SMUX_RECVBUFFER         = 0x400000
	SMUX_STREAMBUFFER       = 0x10000
)

var (
	TIMEOUT = 5000

	PROTOCOL = "TCP"

	// enable log output
	VERBOSE = false

	// logic optimization, changed in v0.1.1
	FORWARD_WITHOUT_DEC = false
)
