package types

import (
	"net/http"
	"strings"

	"github.com/mz-eco/memoir"
)

type HttpHeader http.Header

func (m HttpHeader) Component() memoir.Component {

	var (
		ul = memoir.NewNameValueList("Header")
	)

	for name, value := range m {
		ul.Append(
			name,
			strings.Join(value, ";"),
		)
	}

	return ul
}
