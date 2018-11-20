package httpproxy

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gobwas/glob"
)

type Group struct {
	glob    glob.Glob
	handler Handler
}

type Groups struct {
	handlers []*Group
}

var (
	defaultHandler = &emptyHandler{}
)

func (m *Groups) Match(url *url.URL) bool {

	if len(url.Host) == 0 {
		return false
	}

	var (
		host = strings.Split(url.Host, ":")[0]
	)

	for _, g := range m.handlers {
		fmt.Println(host, g.glob.Match(host))
		if g.glob.Match(host) {
			return true
		}
	}

	return false
}

func (m *Groups) GetHandler(url *url.URL) Handler {

	for _, g := range m.handlers {

		if g.glob.Match(url.Host) {
			return g.handler
		}
	}

	return defaultHandler
}

func NewGroups() *Groups {
	return &Groups{
		handlers: make([]*Group, 0),
	}
}
