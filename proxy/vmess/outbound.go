package vmess

// OutboundHandler ...
type OutboundHandler struct {
	users       []*User
	sessionList []*ClientSession
}

// NewOutboundHandler 构造函数
func NewOutboundHandler() *OutboundHandler {
	return &OutboundHandler{
		users:       make([]*User, 0, 256),
		sessionList: make([]*ClientSession, 0, 256),
	}
}

// AddUser 添加新用户
func (handler *OutboundHandler) AddUser(user *User) {
	handler.users = append(handler.users, user)
}
