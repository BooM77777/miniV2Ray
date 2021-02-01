package network

import (
	"net"
)

// Connection 对网络连接的封装，可以之后添加一些功能
type Connection struct {
	net.Conn
}
