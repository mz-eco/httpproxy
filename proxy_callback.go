package httpproxy

import "github.com/mz-eco/httpproxy/types"

type Handler interface {
	Error(ctx *types.Translate, err error)
	OnRequest(ctx *types.Translate) error
	OnResponse(ctx *types.Translate, err error) error
	Done(ctx *types.Translate)
}
