package httpproxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/mz-eco/httpproxy/utils"

	"context"
)

type proxyHandler struct {
	Schema string
	Proxy  *Proxy
}

func (m *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m.Proxy.onProxyHandler(
		m.Schema,
		w,
		r)
}

type Hosts []string

type Proxy struct {
	https *http.Server

	tls         chan net.Conn
	tlsListener *TunnelListener

	httpProxyHandler  *proxyHandler
	httpsProxyHandler *proxyHandler
	fileHandler       http.Handler

	//TSL
	certFile string
	keyFile  string

	c *http.Client

	source *Source

	//hooks
	hooks *Hooks

	addr    net.Addr
	address string

	running bool
	err     error
}

func (m *Proxy) IsRunning() bool {
	return m.running
}

func (m *Proxy) GetTranslates() int {

	return m.source.GetSize()
}

func (m *Proxy) copyHeaders(source, target http.Header) {

	for key, values := range source {
		for _, value := range values {
			target.Add(key, value)
		}
	}
}

func (m *Proxy) onProxyHandler(schema string, w http.ResponseWriter, r *http.Request) {

	var (
		url = &url.URL{
			Scheme:     schema,
			Opaque:     r.URL.Opaque,
			Host:       r.Host,
			Path:       r.URL.Path,
			RawPath:    r.URL.RawPath,
			ForceQuery: r.URL.ForceQuery,
			RawQuery:   r.URL.RawQuery,
			Fragment:   r.URL.Fragment,
		}
	)

	m.do(
		url,
		w,
		r)
}

func (m *Proxy) ServeProxy(ack http.ResponseWriter, ask *http.Request) {

	if ask.Method == "CONNECT" {

		hijack, ok := ack.(http.Hijacker)

		if !ok {
			fmt.Println("hijack not support for %s fail.", ask.Host)
		}

		ack.WriteHeader(http.StatusOK)

		conn, _, err := hijack.Hijack()

		if err != nil {
			fmt.Println("hijack error", ask.Host)
		}

		m.tls <- conn

	} else {
		m.httpProxyHandler.ServeHTTP(ack, ask)
	}
}

func (m *Proxy) serveTranslate(w http.ResponseWriter, r *http.Request) {

	if r.Method == "CONNECT" {

		hijack, ok := w.(http.Hijacker)

		if !ok {
			fmt.Println("hijack not support for %s fail.", r.Host)
		}

		w.WriteHeader(http.StatusOK)

		conn, _, err := hijack.Hijack()

		if err != nil {
			fmt.Println("hijack error", r.Host)
		}

		_, err = NewTLSTranslate(r.Host, conn)

	} else {

		ask, err := http.NewRequest(
			r.Method,
			r.URL.String(),
			r.Body)

		if err != nil {
			fmt.Println(err)
		}

		utils.CopyHeaders(r.Header, ask.Header)

		response, err := m.c.Do(ask)

		if err != nil {
			fmt.Println(err)
			return
		}

		utils.CopyHeaders(response.Header, w.Header())
		w.WriteHeader(response.StatusCode)
		_, err = io.Copy(w, response.Body)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (m *Proxy) ServeHTTP(ack http.ResponseWriter, ask *http.Request) {

	if ask.URL.Host == "proxy.xz" && ask.URL.Path == "/proxy.crt" {
		m.fileHandler.ServeHTTP(ack, ask)
	}

	if m.hooks.Match(ask.URL) {
		m.ServeProxy(ack, ask)
	} else {
		m.serveTranslate(ack, ask)
	}
}

func (m *Proxy) Server() error {

	m.fileHandler = http.FileServer(
		http.Dir("./"))

	m.https = &http.Server{
		Handler: m.httpsProxyHandler,
	}

	go func() {
		m.err = m.https.ServeTLS(m.tlsListener, m.certFile, m.keyFile)
	}()

	l, err := net.Listen("tcp", m.address)

	if err != nil {
		return err
	}

	m.addr = l.Addr()
	m.running = true

	srv := &http.Server{
		Handler: m,
	}

	srv.RegisterOnShutdown(func() {
		m.running = false
	})

	go srv.Serve(l)

	return nil

}

func (m *Proxy) GetAddr() net.Addr {
	return m.addr
}

func (m *Proxy) GetError() error {
	return m.err
}
func New(opts ...Opt) (*Proxy, error) {

	px := &Proxy{
		hooks:  NewHookers(),
		tls:    make(chan net.Conn),
		source: NewSource(),
		c: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					var (
						d = net.Dialer{
							Timeout:   30 * time.Second,
							KeepAlive: 10 * time.Second,
						}
					)

					return d.DialContext(ctx, network, addr)
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	for _, opt := range opts {
		err := opt(px)

		if err != nil {
			return nil, err
		}
	}

	px.httpProxyHandler = &proxyHandler{Schema: "http", Proxy: px}
	px.httpsProxyHandler = &proxyHandler{Schema: "https", Proxy: px}
	px.tlsListener = NewTunnelListener(px.tls, nil)

	return px, nil
}
