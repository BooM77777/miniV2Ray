package vmess

import (
	"encoding/binary"
	"net"

	"../../common"
)

// TargetAddress 用于存储消息转发的目的地址
type TargetAddress struct {
	addrType   byte   // 地址类型
	ipaddr     net.IP // ip地址
	domainName string // 域名
	port       uint16 // 端口
}

// EncodeTargetAddress 编码目的地址
func (target *TargetAddress) EncodeTargetAddress() (addrBytes []byte, err error) {

	writeBuf := common.GetWriteBuffer()
	defer common.PutWriteBuffer(writeBuf)

	// 编码端口
	err = binary.Write(writeBuf, binary.BigEndian, target.port)
	if err != nil {
		return
	}

	// 编码地址类型
	err = writeBuf.WriteByte(target.addrType)
	if err != nil {
		return
	}

	// 编码地址
	switch target.addrType {
	case AtypIP4:
		_, err = writeBuf.Write(target.ipaddr)
		break
	case AtypDomain:
		err = writeBuf.WriteByte(byte(len(target.domainName)))
		_, err = writeBuf.Write([]byte(target.domainName))
		break
	case AtypIP6:
		_, err = writeBuf.Write(target.ipaddr)
		break
	}

	return writeBuf.Bytes(), err
}

// DecodeTargetAddress 解码目的地址
func (target *TargetAddress) DecodeTargetAddress(input []byte) (err error) {
	return nil
}
