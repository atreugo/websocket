package websocket

import (
	"os"

	"github.com/fasthttp/websocket"
	"github.com/savsgio/atreugo/v11"
	"github.com/savsgio/go-logger/v2"
	"github.com/valyala/fasthttp"
)

var closeCodes = []int{
	websocket.CloseNormalClosure,
	websocket.CloseGoingAway,
	websocket.CloseNoStatusReceived,
}

func defaultErrorView(ctx *atreugo.RequestCtx, err error, statusCode int) {
	ctx.Error(err.Error(), statusCode)
}

func New(cfg Config) *Upgrader {
	if cfg.Error == nil {
		cfg.Error = defaultErrorView
	}

	if cfg.Logger == nil {
		cfg.Logger = logger.New("atreugo-websocket", logger.ERROR, os.Stderr)
	}

	upgrader := &websocket.FastHTTPUpgrader{
		HandshakeTimeout:  cfg.HandshakeTimeout,
		ReadBufferSize:    cfg.ReadBufferSize,
		WriteBufferSize:   cfg.WriteBufferSize,
		Subprotocols:      cfg.Subprotocols,
		EnableCompression: cfg.EnableCompression,
	}

	upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool {
		origin := string(ctx.Request.Header.Peek("Origin"))

		for _, v := range cfg.AllowedOrigins {
			if v == origin || v == "*" {
				return true
			}
		}

		return false
	}

	upgrader.Error = func(ctx *fasthttp.RequestCtx, status int, reason error) {
		actx := acquireRequestCtx(ctx)
		cfg.Error(actx, reason, status)
		releaseRequestCtx(actx)
	}

	return &Upgrader{upgrader: upgrader, logger: cfg.Logger}
}

func (u *Upgrader) Upgrade(viewFn View) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		ws := acquireConn()

		// Copy user values
		ctx.VisitUserValues(func(key []byte, value interface{}) {
			ws.values.SetBytes(key, value)
		})

		return u.upgrader.Upgrade(ctx.RequestCtx, func(conn *websocket.Conn) {
			// Ensure set the connection
			ws.Conn = conn

			if err := viewFn(ws); err != nil && !websocket.IsCloseError(err, closeCodes...) {
				u.logger.Errorf("Websocket - %v", err)
			}

			ws.Close()
			releaseConn(ws)
		})
	}

}
