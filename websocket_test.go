package websocket

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/atreugo/httptest"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

func Test_New(t *testing.T) {
	cfg := Config{
		AllowedOrigins:    []string{"atreugo.test"},
		HandshakeTimeout:  10,
		ReadBufferSize:    20,
		WriteBufferSize:   30,
		Subprotocols:      []string{"test"},
		EnableCompression: true,
	}

	u := New(cfg)

	if cfg.HandshakeTimeout != u.upgrader.HandshakeTimeout {
		t.Errorf("Upgrader.HandshakeTimeout == %d, want %d", cfg.HandshakeTimeout, u.upgrader.HandshakeTimeout)
	}

	if cfg.ReadBufferSize != u.upgrader.ReadBufferSize {
		t.Errorf("Upgrader.ReadBufferSize == %d, want %d", cfg.ReadBufferSize, u.upgrader.ReadBufferSize)
	}

	if cfg.WriteBufferSize != u.upgrader.WriteBufferSize {
		t.Errorf("Upgrader.WriteBufferSize == %d, want %d", cfg.WriteBufferSize, u.upgrader.WriteBufferSize)
	}

	wantSubprotocols := strings.Join(cfg.Subprotocols, ",")
	subprotocols := strings.Join(u.upgrader.Subprotocols, ",")
	if subprotocols != wantSubprotocols {
		t.Errorf("Upgrader.Subprotocols == %s, want %s", subprotocols, wantSubprotocols)
	}

	if cfg.EnableCompression != u.upgrader.EnableCompression {
		t.Errorf("Upgrader.EnableCompression == %v, want %v", cfg.EnableCompression, u.upgrader.EnableCompression)
	}

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.Set("Origin", cfg.AllowedOrigins[0])

	if !u.upgrader.CheckOrigin(ctx) {
		t.Errorf("Upgrader.CheckOrigin == %v, want %v", false, true)
	}

	ctx.Request.Header.Set("Origin", "forbidden.test")

	if u.upgrader.CheckOrigin(ctx) {
		t.Errorf("Upgrader.CheckOrigin == %v, want %v", true, false)
	}

	cfg.AllowedOrigins = []string{"*"}
	u = New(cfg)

	ctx.Request.Header.Set("Origin", "forbidden.test")

	if !u.upgrader.CheckOrigin(ctx) {
		t.Errorf("Upgrader.CheckOrigin == %v, want %v", false, true)
	}

	wantStatusCode := 500
	err := errors.New("some error")

	u.upgrader.Error(ctx, wantStatusCode, err)

	statusCode := ctx.Response.StatusCode()
	body := string(ctx.Response.Body())

	if statusCode != wantStatusCode {
		t.Errorf("Upgrader.Error status code == %d, want %d", statusCode, wantStatusCode)
	}

	wantBody := err.Error()

	if body != wantBody {
		t.Errorf("Upgrader.Error body == %s, want %s", body, wantBody)
	}
}

func Test_Upgrade(t *testing.T) {
	cfg := Config{
		AllowedOrigins: []string{"atreugo.test"},
	}
	u := New(cfg)

	executed := false
	values := make(map[string]interface{})

	hijackDone := make(chan struct{}, 1)
	wsView := u.Upgrade(func(ws *Conn) error {
		executed = true
		ws.values.Map(values)

		close(hijackDone)

		return nil
	})

	fnView := func(ctx *atreugo.RequestCtx) error {
		ctx.Request.Header.Set("Origin", cfg.AllowedOrigins[0])
		ctx.Request.Header.Set("Connection", "Upgrade")
		ctx.Request.Header.Set("Upgrade", "Websocket")
		ctx.Request.Header.Set("Sec-Websocket-Version", "13")
		ctx.Request.Header.Set("Sec-Websocket-Key", "abcd123fgh")

		ctx.SetUserValue("key", "value")

		return wsView(ctx)
	}

	req := new(fasthttp.Request)
	req.SetRequestURI("/ws")
	req.Header.SetMethod("GET")

	httptest.AssertView(t, req, fnView, func(resp *fasthttp.Response) {
		select {
		case <-hijackDone:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("hijack timeout")
		}

		if !executed {
			t.Error("Websocket view is not executed")
		}

		if len(values) == 0 {
			t.Error("UserValues are not saved")
		}
	})
}
