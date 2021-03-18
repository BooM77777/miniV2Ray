package vmess

import (
	"io"
	"net/url"
	"testing"

	"github.com/BooM77777/miniV2Ray/common"
)

type TestIO struct {
	io.Reader
	io.Writer

	buffer []byte

	readIter int
}

func (io *TestIO) Write(input []byte) (n int, err error) {
	io.buffer = append(io.buffer, input...)
	n, err = len(input), nil
	return
}

func (io *TestIO) Read(output []byte) (n int, err error) {

	copy(output, io.buffer[io.readIter:])
	n, err = len(output), nil
	io.readIter += n
	return
}

func TestVmess(t *testing.T) {

	urlStr := "vmess://598c8d05-6ff1-4ead-83f1-0babcefd75b9@127.0.0.1:9527?alterID=4"
	u := common.Must2(url.Parse(urlStr)).(*url.URL)

	io := &TestIO{
		buffer:   []byte{},
		readIter: 0,
	}

	user := NewUser([]byte(u.User.String()))
	distAddr := CreateTargetAddrByDomainName(AddrTypeDomain, "www.google.com", 443)

	inboundHandler := CreateInboundHandler(io, io)
	inboundHandler.AddUser(user)

	clientSession := createNewClientSession(user, distAddr, 0, 0, 0)
	clientSession.auth(io)
	clientSession.encodeRequestHeader(io)

	inboundHandler.CreateSession()

}
