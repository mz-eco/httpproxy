package httpproxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mz-eco/httpproxy/types"

	"github.com/gin-gonic/gin"
)

type ApiTranslate struct {
	Host       string `json:"host"`
	Path       string `json:"path"`
	Error      bool   `json:"error"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Method     string `json:"method"`
	Message    string `json:"message"`
	Index      int    `json:"index"`
}

type View struct {
	URL string    `json:"url"`
	Ask *ViewBody `json:"ask"`
	Ack *ViewBody `json:"ack"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BodyType int

const (
	BodyJSON  BodyType = 0
	BodyBytes          = 1
)

type DataBody struct {
	BodyType BodyType `json:"body_type"`
	String   string   `json:"string"`
	Bytes    []byte   `json:"bytes"`
}

type ViewBody struct {
	Headers []*NameValue `json:"headers"`
	Body    *DataBody    `json:"body"`
}

func NewDataBody(body []byte) *DataBody {

	if body == nil || len(body) == 0 {
		return &DataBody{
			BodyType: BodyBytes,
		}
	}
	if body[0] == '{' || body[0] == '[' {

		var (
			data = make(map[string]interface{})
		)

		err := json.Unmarshal(body, &data)

		if err != nil {
			return &DataBody{
				BodyType: BodyBytes,
				Bytes:    body,
			}
		}

		bytes, err := json.MarshalIndent(data, "", "    ")

		if err != nil {
			return &DataBody{
				BodyType: BodyBytes,
				Bytes:    body,
			}
		}

		return &DataBody{
			BodyType: BodyJSON,
			String:   string(bytes),
			Bytes:    body,
		}
	}

	return &DataBody{
		BodyType: BodyBytes,
		Bytes:    body,
	}
}

func NewViewBody(headers http.Header) *ViewBody {

	vb := &ViewBody{
		Headers: make([]*NameValue, 0),
	}

	if headers != nil {
		for name, value := range headers {
			vb.Headers = append(vb.Headers, &NameValue{
				Name:  name,
				Value: strings.Join(value, ";"),
			})
		}
	}

	return vb
}

func RunApiServer(addr string, s *Source) error {

	e := gin.Default()

	e.GET("/api/view2/:id", func(context *gin.Context) {
		var (
			id int
		)

		fmt.Sscanf(
			context.Params.ByName("id"),
			"%d",
			&id,
		)

		ctx := s.contexts[id]

		context.JSON(
			http.StatusOK,
			ctx.Document())
	})

	e.GET("/api/translate", func(ctx *gin.Context) {

		var (
			tr = make([]*types.Summary, 0)
		)

		for _, ctx := range s.contexts {

			tr = append(tr, ctx.Summary())

		}

		ctx.JSON(
			http.StatusOK,
			&map[string]interface{}{
				"items": tr,
			},
		)

	})

	return http.ListenAndServe(addr, e)
}
