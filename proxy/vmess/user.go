package vmess

import (
	"crypto/md5"

	"../../common"
)

// User ...
type User struct {
	uuid   []byte
	cmdKey []byte
}

// NewUser 构造函数
func NewUser(uuid []byte) *User {

	cmdKey := common.Must2(common.GetBuffer(32)).([]byte)
	defer common.PutBuffer(cmdKey)

	copy(cmdKey[:16], uuid)
	copy(cmdKey[16:], []byte("c48619fe-8f02-49e0-b9e9-edf763e17e21"))

	h := md5.New()
	h.Write(cmdKey)

	return &User{
		uuid:   uuid,
		cmdKey: h.Sum(nil),
	}
}

// GetUUID 获取用户UUID
func (u *User) GetUUID() []byte {
	return u.uuid[:]
}

// GetCmdKey 获取用户GetCmdKey
func (u *User) GetCmdKey() []byte {
	return u.cmdKey
}
