package vmess

// Session vmess会话
type Session struct {
	user   *User          // 用户
	target *TargetAddress // 对端地址

	iv  [16]byte // 数据加密 IV：随机值；
	key [16]byte // 数据加密 Key：随机值；
	v   byte     // 响应认证 V：随机值；
	opt byte     // 选项
	sec byte     // 加密方式
	cmd byte     // 指令（TCP数据或UDP数据）
}
