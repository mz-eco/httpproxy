package types

import (
	"github.com/mz-eco/ui"
)

type HttpBody struct {
	Body []byte
}

func (m *HttpBody) BodyUI() *ui.DataView {
	return ui.NewDataView("Body", m.Body)
}
