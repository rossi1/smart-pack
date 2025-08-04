package httperr

import (
	"net/http"

	"github.com/rossi1/smart-pack/pkg/server/dto"
	"github.com/sirupsen/logrus"
)

type ErrorMessageBody struct {
	Code     int       `json:"code"`
	Messages []Message `json:"messages"`
}

// Message represents an individual error message
type Message struct {
	FormProperty string `json:"formProperty"`
	Label        string `json:"label"`
}

func NewErrorMessage(label, formProperty string) []Message {
	return []Message{
		{
			FormProperty: formProperty,
			Label:        label,
		},
	}
}

func NewErrorMessageBodyWithMessages(messages []Message) *ErrorMessageBody {
	return &ErrorMessageBody{
		Messages: messages,
	}
}

func InternalError(label, formProperty string, err error, w http.ResponseWriter, r *http.Request) {
	m := NewErrorMessage(label, formProperty)
	httpRespondWithError(err, w, r, http.StatusInternalServerError, m)
}

func UnprocessableEntity(label, formProperty string, err error, w http.ResponseWriter, r *http.Request) {
	m := NewErrorMessage(label, formProperty)
	httpRespondWithError(err, w, r, http.StatusUnprocessableEntity, m)
}

func BadRequest(label, formProperty string, err error, w http.ResponseWriter, r *http.Request) {
	m := NewErrorMessage(label, formProperty)
	httpRespondWithError(err, w, r, http.StatusBadRequest, m)
}

func WithStatus(label, formProperty string, err error, w http.ResponseWriter, r *http.Request, status int) {
	m := NewErrorMessage(label, formProperty)
	httpRespondWithError(err, w, r, status, m)
}

func httpRespondWithError(err error, w http.ResponseWriter, r *http.Request, statusCode int, m []Message) {
	ctx := r.Context()
	logrus.WithContext(ctx).Error(err)
	resp := ErrorMessageBody{
		Code:     statusCode,
		Messages: m,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	dto.Write(w, r, &resp)
}
