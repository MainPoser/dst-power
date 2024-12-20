package options

import (
	"net"
)

type ServerOptions struct {
	// 监听地址
	BindAddress net.IP
	// BindPort 监听端口.
	BindPort     int
	ReadTimeout  int64
	WriteTimeout int64
}
