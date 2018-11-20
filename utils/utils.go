package utils

import "net/http"

func CopyHeaders(src http.Header, target http.Header) {

	for key, items := range src {

		for _, value := range items {
			target.Add(key, value)
		}
	}
}
