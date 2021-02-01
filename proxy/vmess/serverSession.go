package vmess

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"io"
	"log"

	"../../common"
)

// ServerSession ...
type ServerSession struct {
	user   *User          // 用户
	target *TargetAddress // 目标转发地址

	// ?
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

func createNewServerSession(user *User) *ServerSession {
	return &ServerSession{
		user: user,
	}
}

// DecodeResponseHeader 解码请求头部
func (session *ServerSession) DecodeResponseHeader(reader io.Reader) error {
	var paddingLen, addressLen int

	block := common.Must2(aes.NewCipher(session.user.GetCmdKey())).(cipher.Block)
	stream := cipher.NewCFBDecrypter(block, session.cmdIV[:])

	requestHeader := common.GetWriteBuffer()
	defer common.PutWriteBuffer(requestHeader)

	buf := common.Must2(common.GetBuffer(41)).([]byte)
	defer common.PutBuffer(buf)

	common.Must2(io.ReadFull(reader, buf))
	stream.XORKeyStream(buf, buf)
	requestHeader.Write(buf)

	// log.Println("->", buf)

	copy(session.requestIV[:], buf[1:17])
	copy(session.requestKey[:], buf[17:33])

	session.v = buf[33]
	session.opt = buf[34]
	paddingLen = int(buf[35] >> 4)
	session.sec = byte(buf[35] & 0xf)
	// 中间有1字节保留位
	session.cmd = buf[37]

	session.target = &TargetAddress{}
	session.target.DecodePort(buf[38:40])
	switch buf[40] {
	case AddrTypeIPv4:
		addressLen = 4
		break
	case AddrTypeIPv6:
		addressLen = 16
		break
	case AddrTypeDomain:
		lenBytes := common.Must2(common.GetBuffer(1)).([]byte)
		defer common.PutBuffer(lenBytes)
		io.ReadFull(reader, lenBytes)
		stream.XORKeyStream(lenBytes, lenBytes)
		requestHeader.Write(lenBytes)
		addressLen = int(lenBytes[0])
		break
	default:
		return errors.New("undefine address type")
	}

	// 获取剩余字节
	remainBuf := common.Must2(common.GetBuffer(addressLen + paddingLen + 4)).([]byte)
	io.ReadFull(reader, remainBuf)
	stream.XORKeyStream(remainBuf, remainBuf)
	requestHeader.Write(remainBuf)

	// 编码目标地址
	session.target.DecodeAddress(buf[40], remainBuf[:addressLen])

	requestHeaderBytes := requestHeader.Bytes()

	log.Println("->", requestHeaderBytes)
	// 最终校验
	fnv1a := fnv.New32()
	common.Must2(fnv1a.Write(requestHeaderBytes[:len(requestHeaderBytes)-4]))

	actualHash := fnv1a.Sum32()
	expectedHash := binary.BigEndian.Uint32(requestHeaderBytes[len(requestHeaderBytes)-4:])

	if actualHash != expectedHash {
		return errors.New("hash is wrong")
	}

	return nil
}
