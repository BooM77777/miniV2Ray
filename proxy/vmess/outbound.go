package vmess

import (
	"net"
	"net/url"

	"../../common"
)

// OutboundHandler ...
type OutboundHandler struct {
	users       User
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
	session := createNewClientSession(dist)
	// 发送客户端认证消息
	handler.connection.Write(common.Must2(session.auth()).([]byte))
	// 发送请求头部，进行密钥等的协商过程
	handler.connection.Write(common.Must2(session.encodeRequestHeader()).([]byte))
	// 将session加入列表中
	handler.sessionList = append(handler.sessionList, session)
}
