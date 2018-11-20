package httpproxy

import "github.com/mz-eco/httpproxy/types"

type Source struct {
	contexts []*types.Translate
}

func (m *Source) Add(ctx *types.Translate) {
	m.contexts = append(m.contexts, ctx)
}

func NewSource() *Source {
	return &Source{
		contexts: make([]*types.Translate, 0),
	}
}
