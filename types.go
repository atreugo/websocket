package websocket

import (
	"time"

	"github.com/fasthttp/websocket"
	"github.com/savsgio/atreugo/v11"
	"github.com/savsgio/dictpool"
	"github.com/savsgio/go-logger/v2"
)

// Upgrader tool to convert the websocket view to an atreugo view.
type Upgrader struct {
	upgrader *websocket.FastHTTPUpgrader
	logger   *logger.Logger
}

// Config configuration for upgrading an HTTP connection to a WebSocket connection.
type Config struct {
	// Specifies either the allowed origins.
	// The "*" wildcard, allow any origin.
	AllowedOrigins []string

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes in bytes. If a buffer
	// size is zero, then buffers allocated by the HTTP server are used. The
	// I/O buffer sizes do not limit the size of the messages that can be sent
	// or received.
	ReadBufferSize, WriteBufferSize int

	// Subprotocols specifies the server's supported protocols in order of
	// preference. If this field is not nil, then the Upgrade method negotiates a
	// subprotocol by selecting the first match in this list with a protocol
	// requested by the client. If there's no match, then no protocol is
	// negotiated (the Sec-Websocket-Protocol header is not included in the
	// handshake response).
	Subprotocols []string

	// EnableCompression specify if the server should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not
	// guarantee that compression will be supported. Currently only "no context
	// takeover" modes are supported.
	EnableCompression bool

	// Error specifies the function for generating HTTP error responses. If Error
	// is nil, then http.Error is used to generate the HTTP response.
	Error atreugo.ErrorView

	// Logger is used for logging formatted messages.
	Logger *logger.Logger
}

// Conn represents a WebSocket connection.
type Conn struct {
	values *dictpool.Dict

	*websocket.Conn
}

// View must process incoming websocket connections.
type View func(ws *Conn) error
