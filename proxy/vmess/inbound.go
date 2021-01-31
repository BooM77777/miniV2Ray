package vmess

// InboundHandler 服务器端运行的入站处理器
type InboundHandler struct {
	userAuthInfoHash map[[16]byte]bool
	session          []*ServerSession
}

// CreaterInboundHandler 构造函数
func CreaterInboundHandler() InboundHandler {
	return InboundHandler{
		userAuthInfoHash: map[[16]byte]bool{},
		session:          []*ServerSession{},
	}
}

// CreateSession ...
func (handler *InboundHandler) CreateSession() {

}
