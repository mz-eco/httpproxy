package types

import (
	"time"

	"github.com/mz-eco/ui"
)

type Summary struct {
	Method     string
	Url        string
	Host       string
	Path       string
	Status     string
	StatusCode int
	Error      bool
	Message    string
	CreateTime time.Time
	TimeUsed   time.Duration
}

func (m *Summary) Component() ui.Component {
	return ui.NewKeyValues(
		"Summary",
		ui.KeyValue{
			"Method":     m.Method,
			"Url":        m.Url,
			"Host":       m.Host,
			"Status":     m.Status,
			"StatusCode": m.StatusCode,
			"Error":      m.Error,
			"Message":    m.Message,
			"CreateTime": m.CreateTime,
			"TimeUsed":   m.TimeUsed,
		})
}
