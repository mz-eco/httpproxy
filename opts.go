package httpproxy

type Opt func(p *Proxy) error

func WithTSL(cert, key string) Opt {
	return func(p *Proxy) error {
		p.certFile = cert
		p.keyFile = key
		return nil
	}
}
