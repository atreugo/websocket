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

// UserValue returns the value stored via SetUserValue* under the given key.
func (ws *Conn) UserValue(key string) interface{} {
	return ws.values.Get(key)
}

// UserValueBytes returns the value stored via SetUserValue* under the given key.
func (ws *Conn) UserValueBytes(key []byte) interface{} {
	return ws.values.GetBytes(key)
}

// SetUserValue stores the given value (arbitrary object)
// under the given key.
//
// The value stored may be obtained by UserValue*.
func (ws *Conn) SetUserValue(key string, value interface{}) {
	ws.values.Set(key, value)
}

// SetUserValueBytes stores the given value (arbitrary object)
// under the given key.
//
// The value stored may be obtained by UserValue*.
func (ws *Conn) SetUserValueBytes(key []byte, value interface{}) {
	ws.values.SetBytes(key, value)
}
