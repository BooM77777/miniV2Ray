package vmess

import "net"

// TargetAddress 用于存储消息转发的目的地址
type TargetAddress struct {
	ipaddr     net.IP // ip地址
	domainName string // 域名
	port       uint16 // 端口
}
