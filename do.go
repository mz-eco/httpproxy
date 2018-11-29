package httpproxy

import (
	"net/url"

	"github.com/mz-eco/http"

	H "net/http"
)

func (m *Proxy) do(url *url.URL, w H.ResponseWriter, r *H.Request) {

	var (
		x = m.hooks.GetGroup(url)
	)

	m.source.Add(
		http.Do(
			x,
			http.Clone(r).SetURL(url.String()),
		),
	)
}
