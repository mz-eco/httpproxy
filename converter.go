package httpproxy

import "net/http"

type Converter interface {
	ConvertHeader(header http.Header) http.Header
	ConvertBody(bytes []byte) []byte
}
