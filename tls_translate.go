package httpproxy

import (
	"net"
)

type TLSTranslate struct {
	local  net.Conn
	remote net.Conn
}

func (m *TLSTranslate) io(src net.Conn, target net.Conn) {

	var (
		buf = make([]byte, 10240)
	)

	for {
		n, err := src.Read(buf)

		if err != nil {
			return
		}

		_, err = target.Write(buf[0:n])

		if err != nil {
			return
		}
	}
}

func NewTLSTranslate(addr string, local net.Conn) (*TLSTranslate, error) {

	var (
		t = &TLSTranslate{
			local: local,
		}

		err error
	)

	t.remote, err = net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	go t.io(t.remote, t.local)
	go t.io(t.local, t.remote)

	return t, nil

}
