package types

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mz-eco/ui"
)

type Translate struct {
	CreateTime time.Time
	Ask        *Request
	Ack        *Response
	Error      error
	URL        *url.URL
}

func NewContext() *Translate {

	return &Translate{
		CreateTime: time.Now(),
		Ask: &Request{
			Value: &http.Request{},
		},
		Ack: &Response{
			Value: &http.Response{},
		},
	}
}

func (m *Translate) Message() string {

	if m.Error != nil {
		return m.Error.Error()
	}

	return ""
}

func (m *Translate) Summary() *Summary {
	return &Summary{
		Method:     m.Ask.Value.Method,
		Url:        fmt.Sprintf("%s://%s%s", m.URL.Scheme, m.URL.Host, m.URL.Path),
		Host:       m.URL.Host,
		Status:     http.StatusText(m.Ack.Value.StatusCode),
		StatusCode: m.Ack.Value.StatusCode,
		Error:      m.Error != nil,
		Message:    m.Message(),
		CreateTime: m.CreateTime,
		TimeUsed:   10 * time.Second,
	}
}

func (m *Translate) Document() *ui.Document {

	return ui.NewDocument(
		ui.DocHtmlTranslate,
		"HttpTranslate",
		m.Summary(),
		ui.NewLabel("Value",
			HttpQuery(m.URL.Query()),
			HttpHeader(m.Ask.Value.Header),
			ui.NewDataView(
				"Body",
				m.Ask.Body),
		),
		ui.NewLabel("Value",
			HttpHeader(m.Ack.Value.Header),
			ui.NewDataView(
				"Body",
				m.Ack.Body),
		),
	)

}