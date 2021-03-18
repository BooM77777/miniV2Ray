package vmess

import (
	"io"
	"net/url"

	"github.com/BooM77777/miniV2Ray/common"
)

// OutboundHandler ...
type OutboundHandler struct {
	user        *User
	sessionList []*ClientSession

	// connection *net.Conn
	dataReader io.Reader
	dataWriter io.Writer
}

// CreateOutboundHandlerByURL 根据URL创建出站处理器
func CreateOutboundHandlerByURL(urlStr string, dataReader io.Reader, dataWriter io.Writer) *OutboundHandler {

	u := common.Must2(url.Parse(urlStr)).(*url.URL)

	outboundHandler := &OutboundHandler{
		user:        NewUser([]byte(u.User.String())),
		sessionList: make([]*ClientSession, 0, 256),

		// connection:  connection,
		dataReader: dataReader,
		dataWriter: dataWriter,
	}

	return outboundHandler
}

// CreateSession 创建新的出站会话
func (handler *OutboundHandler) CreateSession(dist *TargetAddress) {
	// 创建新的会话
	session := createNewClientSession(handler.user, dist, 0, 0, 0)
	// 发送客户端认证消息
	session.auth(handler.dataWriter)
	// 发送请求头部，进行密钥等的协商过程
	session.encodeRequestHeader(handler.dataWriter)
	// 将session加入列表中
	handler.sessionList = append(handler.sessionList, session)
}
