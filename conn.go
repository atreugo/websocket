package websocket

import (
	"sync"

	"github.com/savsgio/gotils/strconv"
)

var connPool = sync.Pool{
	New: func() interface{} {
		ws := new(Conn)
		ws.values = make(map[string]interface{})

		return ws
	},
}

func acquireConn() *Conn {
	return connPool.Get().(*Conn) // nolint:forcetypeassert
}

func releaseConn(ws *Conn) {
	ws.reset()
	connPool.Put(ws)
}

func (ws *Conn) reset() {
	for k := range ws.values {
		delete(ws.values, k)
	}

	ws.Conn = nil
}

// UserValue returns the value stored via SetUserValue* under the given key.
func (ws *Conn) UserValue(key string) interface{} {
	return ws.values[key]
}

// UserValueBytes returns the value stored via SetUserValue* under the given key.
func (ws *Conn) UserValueBytes(key []byte) interface{} {
	return ws.UserValue(strconv.B2S(key))
}

// SetUserValue stores the given value (arbitrary object)
// under the given key.
//
// The value stored may be obtained by UserValue*.
func (ws *Conn) SetUserValue(key string, value interface{}) {
	ws.values[key] = value
}

// SetUserValueBytes stores the given value (arbitrary object)
// under the given key.
//
// The value stored may be obtained by UserValue*.
func (ws *Conn) SetUserValueBytes(key []byte, value interface{}) {
	ws.SetUserValue(string(key), value)
}
