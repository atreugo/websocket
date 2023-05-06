package websocket

import (
	"testing"

	"github.com/fasthttp/websocket"
)

func Test_aquireConn(t *testing.T) {
	if conn := acquireConn(); conn == nil {
		t.Error("Conn is nil")
	}
}

func Test_releaseConn(t *testing.T) {
	wConn := new(websocket.Conn)

	conn := acquireConn()
	conn.Conn = wConn
	conn.SetUserValue("key", "value")

	releaseConn(conn)

	if conn.Conn != nil {
		t.Error("Conn is not nil")
	}

	if len(conn.values) > 0 {
		t.Error("Conn.values has not been reset")
	}
}

func Test_SetAndGetUserValue(t *testing.T) {
	wConn := new(websocket.Conn)

	conn := acquireConn()
	conn.Conn = wConn

	conn.SetUserValue("key", "value")

	if val := conn.UserValue("key"); val == nil {
		t.Error("UserValue is not saved")
	}
}

func Test_SetAndGetUserValueBytes(t *testing.T) {
	wConn := new(websocket.Conn)

	conn := acquireConn()
	conn.Conn = wConn

	conn.SetUserValueBytes([]byte("key"), "value")

	if val := conn.UserValueBytes([]byte("key")); val == nil {
		t.Error("UserValue is not saved")
	}
}
