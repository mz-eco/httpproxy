package types

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mz-eco/httpproxy/utils"
)

type Request struct {
	HttpBody
	Value *http.Request
}

func (m *Request) GetRequest(url *url.URL) (*http.Request, error) {

	ask, err := http.NewRequest(
		m.Value.Method,
		url.String(),
		bytes.NewBuffer(m.Body))

	if err != nil {
		return nil, err
	}

	utils.CopyHeaders(m.Value.Header, ask.Header)

	return ask, nil
}

func NewHttpRequest(r *http.Request) (*Request, error) {

	var (
		ask = &Request{
			Value: r,
			HttpBody: HttpBody{
				Body: make([]byte, 0),
			},
		}

		err error
	)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	ask.Body = body

	return ask, nil

}
