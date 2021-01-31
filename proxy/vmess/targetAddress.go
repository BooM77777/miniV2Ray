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
	case AddrTypeIPv4:
		_, err = writeBuf.Write(target.ipaddr)
		break
	case AddrTypeIPv6:
		_, err = writeBuf.Write(target.ipaddr)
		break
	case AddrTypeDomain:
		err = writeBuf.WriteByte(byte(len(target.domainName)))
		_, err = writeBuf.Write([]byte(target.domainName))
		break
	}

	return writeBuf.Bytes(), err
}

// DecodePort 解码port
func (target *TargetAddress) DecodePort(portBytes []byte) {
	target.port = binary.BigEndian.Uint16(portBytes)
}

// DecodeAddress 解码目的IP或域名
func (target *TargetAddress) DecodeAddress(addrType byte, addrBytes []byte) {

	switch addrType {
	case AddrTypeIPv4:
		target.ipaddr = make(net.IP, 4)
		copy(target.ipaddr, addrBytes)
		break
	case AddrTypeIPv6:
		target.ipaddr = make(net.IP, 16)
		copy(target.ipaddr, addrBytes)
		break
	case AddrTypeDomain:
		target.domainName = string(addrBytes)
		break
	}
}

// DecodeTargetAddress 解码目的地址
func (target *TargetAddress) DecodeTargetAddress(input []byte) (err error) {
	return nil
}
