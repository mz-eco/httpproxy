package httpproxy

import "github.com/mz-eco/httpproxy/types"

type emptyHandler struct {
}

func (emptyHandler) Error(ctx *types.Translate, err error) {

}

func (emptyHandler) OnRequest(ctx *types.Translate) error {
	return nil
}

func (emptyHandler) OnResponse(ctx *types.Translate, err error) error {
	return nil
}

func (emptyHandler) Done(ctx *types.Translate) {
}
