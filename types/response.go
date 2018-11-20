package types

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mz-eco/httpproxy/utils"
)

type Response struct {
	HttpBody
	Value *http.Response
}

func NewResponse(response *http.Response) (*Response, error) {

	var (
		ack = &Response{
			HttpBody{
				Body: make([]byte, 0),
			},
			response,
		}

		err error
	)

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	ack.Body = body

	return ack, nil

}

func (m *Response) Write(w http.ResponseWriter) error {

	utils.CopyHeaders(m.Value.Header, w.Header())
	w.WriteHeader(m.Value.StatusCode)

	_, err := io.Copy(w, bytes.NewBuffer(m.Body))

	return err
}
