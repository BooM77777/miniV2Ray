package vmess

import (
	"net"
)

// OutboundConn ...
type OutboundConn struct {
	user *User
	conn net.Conn

	iv  [16]byte // 数据加密 IV：随机值；
	key [16]byte // 数据加密 Key：随机值；
	v   byte     // 响应认证 V：随机值；
	opt byte     // 选项
	sec byte     // 加密方式
	cmd byte     // 指令（TCP数据或UDP数据）

	target *TargetAddress // 目的地址
}
