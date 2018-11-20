package types

import (
	"net/url"
	"strings"

	"github.com/mz-eco/memoir"
)

type HttpQuery url.Values

func (m HttpQuery) Component() memoir.Component {

	var (
		ul = memoir.NewNameValueList("Query")
	)

	for name, value := range m {
		ul.Append(
			name,
			strings.Join(value, ","),
		)
	}

	return ul

}
