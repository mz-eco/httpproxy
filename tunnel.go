package httpproxy

import (
	"io"
	"net"
)

type TunnelListener struct {
	net.Listener
	tls   chan net.Conn
	close chan struct{}
	addr  net.Addr
}

func NewTunnelListener(tls chan net.Conn, addr net.Addr) *TunnelListener {

	return &TunnelListener{
		tls:   tls,
		close: make(chan struct{}),
		addr:  addr,
	}
}

func (m *TunnelListener) Accept() (net.Conn, error) {

	for {
		select {
		case <-m.close:
			return nil, io.EOF
		case tls := <-m.tls:
			return tls, nil
		}
	}
}

func (m *TunnelListener) Close() error {
	close(m.close)
	return nil
}

func (m *TunnelListener) Addr() net.Addr {
	return m.addr
}
