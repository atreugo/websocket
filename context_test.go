package websocket

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func Test_acquireRequestCtx(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	actx := acquireRequestCtx(ctx)

	if actx == nil {
		t.Error("atreugo.RequestCtx is nil")
	}
}

func Test_releaseRequestCtx(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	actx := acquireRequestCtx(ctx)
	releaseRequestCtx(actx)

	if actx.RequestCtx != nil {
		t.Error("Internal RequestCtx is not nil")
	}
}
