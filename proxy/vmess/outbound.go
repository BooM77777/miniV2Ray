package vmess

import (
	"net"
	"net/url"
)

// OutboundHandler ...
type OutboundHandler struct {
	user        *User
	sessionList []*ClientSession

	connection net.Conn
}

// CreateOutboundHandlerByURL 根据URL创建出站处理器
func CreateOutboundHandlerByURL(url url.URL) *OutboundHandler {
	return nil
}

// CreateSession 创建新的出站会话
func (handler *OutboundHandler) CreateSession(dist *TargetAddress) {
	// 创建新的会话
	session := createNewClientSession(handler.user, dist, 0, 0, 0)
	// 发送客户端认证消息
	session.auth(handler.connection)
	// 发送请求头部，进行密钥等的协商过程
	session.encodeRequestHeader(handler.connection)
	// 将session加入列表中
	handler.sessionList = append(handler.sessionList, session)
}
