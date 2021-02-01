package vmess

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"
)

// InboundHandler 服务器端运行的入站处理器
type InboundHandler struct {
	userList         []*User
	userAuthInfoHash map[[16]byte]*UserAtTime
	sessionList      []*ServerSession

	ticker *time.Ticker

	dataReader io.Reader
	dataWriter io.Writer
	// connection net.Conn

	// 读写锁
	mutex sync.RWMutex
}

// CreateInboundHandler 构造函数
func CreateInboundHandler(dataReader io.Reader, dataWriter io.Writer) *InboundHandler {

	handler := &InboundHandler{}

	handler.userList = make([]*User, 0, 256)
	handler.userAuthInfoHash = map[[16]byte]*UserAtTime{}
	handler.sessionList = []*ServerSession{}

	handler.dataReader = dataReader
	handler.dataWriter = dataWriter

	handler.ticker = time.NewTicker(60)

	// 创建协程已控制退出和刷新
	go func() {
		for {
			select {
			// 时间到了就刷新
			case <-handler.ticker.C:
				handler.refresh()
			}
		}
	}()

	return handler
}

// AddUser 添加用户
func (handler *InboundHandler) AddUser(user *User) {
	handler.userList = append(handler.userList, user)
	handler.refresh()
}

func (handler *InboundHandler) refresh() {

	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	for k := range handler.userAuthInfoHash {
		delete(handler.userAuthInfoHash, k)
	}

	// 读取当前时间，以秒为单位
	second := int(time.Now().UTC().Unix())
	for _, user := range handler.userList {
		h := hmac.New(md5.New, user.GetUUID())
		for i := 0; i <= 10; i++ {
			h.Reset()
			var timeStamp [8]byte
			var authHash [16]byte
			binary.BigEndian.PutUint64(timeStamp[:], uint64(second+i-5))
			h.Write(timeStamp[:])
			copy(authHash[:], h.Sum(nil))
			handler.userAuthInfoHash[authHash] = NewUserAtTime(user, timeStamp)
		}
	}
}

// CreateSession ...
func (handler *InboundHandler) CreateSession() error {
	handler.mutex.RLock()
	handler.mutex.RUnlock()
	var authInfo [16]byte
	handler.dataReader.Read(authInfo[:])
	if userAtTime, ok := handler.userAuthInfoHash[authInfo]; ok {
		session := createNewServerSession(userAtTime.user)
		session.cmdIV = getCmdIVByTimestamp(userAtTime.timeStamp)
		session.DecodeResponseHeader(handler.dataReader)
		handler.sessionList = append(handler.sessionList, session)

		return nil
	}
	return errors.New("异常客户端")
}
