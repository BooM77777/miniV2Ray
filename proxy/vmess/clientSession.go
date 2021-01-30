package vmess

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"hash/fnv"
	"math/rand"
	"time"

	"../../common"
)

// ClientSession ...
type ClientSession struct {
	user   *User          // 用户
	target *TargetAddress // 对端地址

	iv  [16]byte // 数据加密 IV：随机值；
	key [16]byte // 数据加密 Key：随机值；
	v   byte     // 响应认证 V：随机值；
	opt byte     // 选项
	sec byte     // 加密方式
	cmd byte     // 指令（TCP数据或UDP数据）
}

// Auth 发送用于验证客户端真实性的16B消息，只能由客户端发送
func (s *ClientSession) Auth() (authBytes []byte, err error) {

	var ts []byte

	// UTC 时间，精确到秒，取值为当前时间的前后 30 秒随机值(8 字节, Big Endian)
	ts, err = common.GetBuffer(8)
	defer common.PutBuffer(ts)
	if err != nil {
		return
	}
	binary.BigEndian.PutUint64(ts, uint64(time.Now().UTC().Unix()))

	// HMAC(H, K, M)
	h := hmac.New(md5.New, s.user.GetUUID())
	h.Write(ts)

	authBytes = h.Sum(nil)

	return
}

// EncodeRequestHeader 编码一个Session的请求头部
func (s *ClientSession) EncodeRequestHeader() (retBytes []byte, err error) {

	var buf []byte

	writeBuf := common.GetWriteBuffer()
	defer common.PutWriteBuffer(writeBuf)

	buf, err = common.GetBuffer(38)
	defer common.PutBuffer(buf)
	if err != nil {
		return nil, err
	}

	buf[0] = 1
	copy(buf[1:17], s.iv[:])
	copy(buf[17:33], s.key[:])
	buf[33] = s.v
	buf[34] = s.opt

	paddingLen := rand.Intn(16)
	buf[35] = byte(paddingLen<<4) | s.sec
	buf[36] = 0
	buf[37] = CmdTCP

	_, err = writeBuf.Write(buf) // 写入前37个定长的子节

	// 写入目的地址
	var targetBytes []byte
	targetBytes, err = s.target.EncodeTargetAddress()
	if err != nil {
		return nil, err
	}
	writeBuf.Write(targetBytes)

	// 写入随机填充
	if paddingLen > 0 {
		var paddingBytes []byte
		paddingBytes, err = common.GetBuffer(paddingLen)
		defer common.PutBuffer(paddingBytes)
		rand.Read(paddingBytes)
		writeBuf.Write(paddingBytes)
	}

	// 写入4子节校验值
	fnv1a := fnv.New32()
	_, err = fnv1a.Write(writeBuf.Bytes())
	if err != nil {
		return nil, err
	}
	writeBuf.Write(fnv1a.Sum(nil))

	retBytes = writeBuf.Bytes()

	return
}

// DecodeRequestHeader 解码一个请求头部
func (s *ClientSession) DecodeRequestHeader(requestHeader []byte) {

}
