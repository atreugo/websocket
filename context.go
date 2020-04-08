package websocket

import (
	"sync"

	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

var requestCtxPool = sync.Pool{
	New: func() interface{} {
		return new(atreugo.RequestCtx)
	},
}

func acquireRequestCtx(ctx *fasthttp.RequestCtx) *atreugo.RequestCtx {
	actx := requestCtxPool.Get().(*atreugo.RequestCtx)
	actx.RequestCtx = ctx

	return actx
}

func releaseRequestCtx(actx *atreugo.RequestCtx) {
	actx.RequestCtx = nil

	requestCtxPool.Put(actx)
}
