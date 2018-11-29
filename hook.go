package httpproxy

import (
	"net/url"
	"strings"

	"github.com/mz-eco/http"

	"github.com/gobwas/glob"
)

type Hooks struct {
	hooks map[glob.Glob]*http.Group
}

func (m Hooks) Add(glob glob.Glob, x *http.Group) {

	m.hooks[glob] = x
}

func (m *Hooks) Match(url *url.URL) bool {

	if len(url.Host) == 0 {
		return false
	}

	var (
		host = strings.Split(url.Host, ":")[0]
	)

	for glob := range m.hooks {
		if glob.Match(host) {
			return true
		}
	}

	return false
}

func (m *Hooks) GetGroup(url *url.URL) *http.Group {

	for glob, x := range m.hooks {

		if glob.Match(url.Host) {
			return x
		}
	}

	panic("logic error.")
}

func NewHookers() *Hooks {
	return &Hooks{
		hooks: make(map[glob.Glob]*http.Group),
	}
}

func (m *Proxy) Hook(host string, x *http.Group) {
	m.hooks.Add(glob.MustCompile(host), x)
}
