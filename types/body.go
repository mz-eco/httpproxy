package types

import (
	"github.com/mz-eco/memoir"
)

type HttpBody struct {
	Body []byte
}

func (m *HttpBody) BodyUI() *memoir.DataView {
	return memoir.NewDataView("Body", m.Body)
}
