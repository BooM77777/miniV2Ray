package vmess

// ServerSession ...
type ServerSession struct {
	user   *User          // 用户
	target *TargetAddress // 目标转发地址

	timeStamp [8]byte  // 时间戳
	cmdIV     [16]byte //

	requestIV  [16]byte // 加密用初始化向量
	requestKey [16]byte // 加密用密钥

	respondIV  [16]byte // 解密用初始化向量
	respondKey [16]byte // 解密用密钥

	v   byte // 响应认证 V：随机值；
	opt byte // 选项
	sec byte // 加密方式
	cmd byte // 指令（TCP数据或UDP数据）
}

func decodeAuth(authBytes []byte) {

}

func decodeRequest() {

}
