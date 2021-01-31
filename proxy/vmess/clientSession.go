package vmess

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"hash/fnv"
	"io"
	"math/rand"
	"net"
	"time"

	"../../common"
)

// ClientSession ...
type ClientSession struct {
	user   *User          // 用户
	target *TargetAddress // 对端地址

	timeStamp [8]byte // 时间戳
	cmdIV     [16]byte

	requestIV  [16]byte // 加密用初始化向量
	requestKey [16]byte // 加密用密钥

	respondIV  [16]byte // 解密用初始化向量
	respondKey [16]byte // 解密用密钥

	v   byte // 响应认证 V：随机值；
	opt byte // 选项
	sec byte // 加密方式
	cmd byte // 指令（TCP数据或UDP数据）
}

func createNewClientSession(user *User, distAddr *TargetAddress, opt, sec, cmd byte) *ClientSession {

	session := &ClientSession{}

	// 填写user和实际服务器地址
	session.user, session.target = user, distAddr

	// 创建随机的IV与密钥
	randBytes := common.Must2(common.GetBuffer(32)).([]byte)
	defer common.PutBuffer(randBytes)

	rand.Read(randBytes)
	copy(session.requestIV[:], randBytes[:16])
	copy(session.requestKey[:], randBytes[16:])

	session.respondIV = md5.Sum(session.requestIV[:])
	session.respondKey = md5.Sum(session.requestKey[:])

	// 给一些其他的东西赋值
	session.v = byte(rand.Intn(256))

	session.cmd = cmd
	session.opt = opt
	session.sec = sec

	return session
}

// auth 发送用于验证客户端真实性的16B消息，只能由客户端发送
func (s *ClientSession) auth(writer io.Writer) {

	// UTC 时间，精确到秒，取值为当前时间的前后 30 秒随机值(8 字节, Big Endian)
	binary.BigEndian.PutUint64(s.timeStamp[:], uint64(time.Now().UTC().Unix()))

	// cmdIV = md5(ts+ts+ts+ts)
	ivsrc := common.Must2(common.GetBuffer(8 * 4)).([]byte)
	defer common.PutBuffer(ivsrc)
	s.cmdIV = md5.Sum(ivsrc)

	// HMAC(H, K, M)
	h := hmac.New(md5.New, s.user.GetUUID())
	h.Write(s.timeStamp[:])

	common.Must2(writer.Write(h.Sum(nil)))
}

// encodeRequestHeader 编码一个Session的请求头部
func (s *ClientSession) encodeRequestHeader(writer io.Writer) {

	var buf, retBytes []byte

	writeBuf := common.GetWriteBuffer()
	defer common.PutWriteBuffer(writeBuf)

	buf = common.Must2(common.GetBuffer(38)).([]byte)
	defer common.PutBuffer(buf)

	buf[0] = 1
	copy(buf[1:17], s.requestIV[:])
	copy(buf[17:33], s.requestKey[:])
	buf[33] = s.v
	buf[34] = s.opt

	paddingLen := rand.Intn(16)
	buf[35] = byte(paddingLen<<4) | s.sec
	buf[36] = 0
	buf[37] = CmdTCP

	common.Must2(writeBuf.Write(buf)) // 写入前37个定长的子节

	// 写入目的地址
	var targetBytes []byte
	targetBytes = common.Must2(s.target.EncodeTargetAddress()).([]byte)
	writeBuf.Write(targetBytes)

	// 写入随机填充
	if paddingLen > 0 {
		var paddingBytes []byte
		paddingBytes = common.Must2(common.GetBuffer(paddingLen)).([]byte)
		defer common.PutBuffer(paddingBytes)
		rand.Read(paddingBytes)
		writeBuf.Write(paddingBytes)
	}

	// 写入4子节校验值
	fnv1a := fnv.New32()
	common.Must2(fnv1a.Write(writeBuf.Bytes()))
	writeBuf.Write(fnv1a.Sum(nil))

	retBytes = writeBuf.Bytes()

	// AES-128-CFB 加密
	block := common.Must2(aes.NewCipher(s.user.GetCmdKey())).(cipher.Block)
	stream := cipher.NewCFBEncrypter(block, s.cmdIV[:])
	stream.XORKeyStream(retBytes, retBytes)

	common.Must2(writer.Write(retBytes))
}

func (s *ClientSession) decodeResponseHeader(c net.Conn) {

}

// Read 从网络流中读取数据
func Read() {

}

// Write 向网络中写入数据
func Write() {

}
