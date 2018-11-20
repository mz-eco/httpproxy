package types

import (
	"net/http"
	"strings"

	"github.com/mz-eco/ui"
)

type HttpHeader http.Header

func (m HttpHeader) Component() ui.Component {

	var (
		ul = ui.NewNameValueList("Header")
	)

	for name, value := range m {
		ul.Append(
			name,
			strings.Join(value, ";"),
		)
	}

	return ul
}
