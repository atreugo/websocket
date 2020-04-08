package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/savsgio/atreugo/v11"
	"github.com/savsgio/go-logger"
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

	if cfg.Error == nil {
		cfg.Error = defaultErrorView
	}

	upgrader.Error = func(ctx *fasthttp.RequestCtx, status int, reason error) {
		actx := acquireRequestCtx(ctx)
		cfg.Error(actx, reason, status)
		releaseRequestCtx(actx)
	}

	return &Upgrader{FastHTTPUpgrader: upgrader}
}

func (u *Upgrader) Upgrade(viewFn View) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		ws := acquireConn()

		// Copy user values
		ctx.VisitUserValues(func(key []byte, value interface{}) {
			ws.values.SetBytes(key, value)
		})

		return u.FastHTTPUpgrader.Upgrade(ctx.RequestCtx, func(conn *websocket.Conn) {
			// Ensure set the connection
			ws.Conn = conn

			if err := viewFn(ws); err != nil && !websocket.IsCloseError(err, closeCodes...) {
				logger.Errorf("Websocket - %v", err)
			}

			releaseConn(ws)
		})
	}

}
