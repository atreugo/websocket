package websocket

import (
	"sync"

	"github.com/savsgio/dictpool"
)

var connPool = &sync.Pool{
	New: func() interface{} {
		ws := new(Conn)
		ws.values = new(dictpool.Dict)

		return ws
	},
}

func acquireConn() *Conn {
	return connPool.Get().(*Conn)
}

func releaseConn(ws *Conn) {
	ws.reset()
	connPool.Put(ws)
}

func (ws *Conn) reset() {
	ws.values.Reset()
	ws.Conn = nil
}

func (ws *Conn) UserValue(key string) interface{} {
	return ws.values.Get(key)
}

func (ws *Conn) SetUserValue(key string, value interface{}) {
	ws.values.Set(key, value)
}
