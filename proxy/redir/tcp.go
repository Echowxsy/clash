package redir

import (
	"net"

	"github.com/Echowxsy/clash/adapters/inbound"
	C "github.com/Echowxsy/clash/constant"
	"github.com/Echowxsy/clash/log"
	"github.com/Echowxsy/clash/tunnel"
)

type RedirListener struct {
	net.Listener
	address string
	closed  bool
}

func NewRedirProxy(addr string) (*RedirListener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	rl := &RedirListener{l, addr, false}

	go func() {
		log.Infoln("Redir proxy listening at: %s", addr)
		for {
			c, err := l.Accept()
			if err != nil {
				if rl.closed {
					break
				}
				continue
			}
			go handleRedir(c)
		}
	}()

	return rl, nil
}

func (l *RedirListener) Close() {
	l.closed = true
	l.Listener.Close()
}

func (l *RedirListener) Address() string {
	return l.address
}

func handleRedir(conn net.Conn) {
	target, err := parserPacket(conn)
	if err != nil {
		conn.Close()
		return
	}
	conn.(*net.TCPConn).SetKeepAlive(true)
	tunnel.Add(inbound.NewSocket(target, conn, C.REDIR, C.TCP))
}
