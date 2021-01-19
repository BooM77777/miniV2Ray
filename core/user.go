package core

// User ...
type User struct {
	uuid   [16]byte
	cmdKey [16]byte
}

// NewUser 构造函数
func NewUser(uuid, cmdKey [16]byte) *User {
	return &User{
		uuid:   uuid,
		cmdKey: cmdKey,
	}
}

// GetUUID 获取用户UUID
func (u *User) GetUUID() [16]byte {
	return u.uuid
}

// GetCmdKey 获取用户GetCmdKey
func (u *User) GetCmdKey() [16]byte {
	return u.cmdKey
}
