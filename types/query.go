package types

import (
	"net/url"
	"strings"

	"github.com/mz-eco/ui"
)

type HttpQuery url.Values

func (m HttpQuery) Component() ui.Component {

	var (
		ul = ui.NewNameValueList("Query")
	)

	for name, value := range m {
		ul.Append(
			name,
			strings.Join(value, ","),
		)
	}

	return ul

}
