package types

import (
	"time"

	"github.com/mz-eco/memoir"
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

func (m *Summary) Component() memoir.Component {
	return memoir.NewKeyValues(
		"Summary",
		memoir.KeyValue{
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
