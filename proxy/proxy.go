package proxy

import (
	"errors"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

//OutboundClient ...TODO
type OutboundClient interface {
	Handshake(underlay net.Conn, target string) (io.ReadWriter, error)
	Process()
}

// ClientCreator is a function to create client.
type ClientCreator func(url *url.URL) (OutboundClient, error)

var (
	clientMap = make(map[string]ClientCreator)
)

//RegisterClient ...
func RegisterClient(name string, c ClientCreator) {
	clientMap[name] = c
}

//ClientFromURL ...
func ClientFromURL(s string) (OutboundClient, error) {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("can not parse client url %s err: %s", s, err)
		return nil, err
	}

	c, ok := clientMap[strings.ToLower(u.Scheme)]
	if ok {
		return c(u)
	}

	return nil, errors.New("unknown client scheme '" + u.Scheme + "'")
}

// InboundServer ...(TODO)
type InboundServer interface {
	Handshake(underlay net.Conn) (io.ReadWriter, error)
	Stop()
}

// ServerCreator is a function to create proxy server
type ServerCreator func(url *url.URL) (InboundServer, error)

var (
	serverMap = make(map[string]ServerCreator)
)

// RegisterServer is used to register a proxy server
func RegisterServer(name string, c ServerCreator) {
	serverMap[name] = c
}

// ServerFromURL calls the registered creator to create proxy servers
// dialer is the default upstream dialer so cannot be nil, we can use Default when calling this function
func ServerFromURL(s string) (InboundServer, error) {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("can not parse server url %s err: %s", s, err)
		return nil, err
	}

	c, ok := serverMap[strings.ToLower(u.Scheme)]
	if ok {
		return c(u)
	}

	return nil, errors.New("unknown server scheme '" + u.Scheme + "'")
}
