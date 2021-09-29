package http

import (
	"net/http"

	json "github.com/json-iterator/go"
	"github.com/julb/go/pkg/util/date"
)

type HttpErrorResponse struct {
	Status    int      `json:"httpStatus"`
	Message   string   `json:"message"`
	DateTime  string   `json:"dateTime"`
	RequestId string   `json:"requestId"`
	Trace     []string `json:"traceId"`
}

// Write a JSON response by marshalling the body as JSON.
func WriteJSON(r *http.Request, w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set(HdrContentType, "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body) // nolint:errcheck
}

// Write a HTTP 404 with the writer with an empty message.
func WriteNotFound(r *http.Request, w http.ResponseWriter) {
	WriteNotFoundM(r, w, "")
}

// Write a HTTP 404 with the writer and include the message.
func WriteNotFoundM(r *http.Request, w http.ResponseWriter, message string) {
	write4xx5xx(r, w, http.StatusNotFound, message)
}

// Write a HTTP 404 with the writer with an empty message.
func WriteServiceUnavailable(r *http.Request, w http.ResponseWriter) {
	WriteServiceUnavailableM(r, w, "")
}

// Write a HTTP 404 with the writer and include the message.
func WriteServiceUnavailableM(r *http.Request, w http.ResponseWriter, message string) {
	write4xx5xx(r, w, http.StatusServiceUnavailable, message)
}

// Write a HTTP 500 with the writer and include the message.
func WriteInternalServerError(r *http.Request, w http.ResponseWriter, message string) {
	write4xx5xx(r, w, http.StatusInternalServerError, message)
}

// Write a 4xx/5xx HTTP response.
func write4xx5xx(r *http.Request, w http.ResponseWriter, status int, message string) {
	WriteJSON(r, w, status, &HttpErrorResponse{
		Status:    status,
		RequestId: w.Header().Get(HdrXRequestId),
		Message:   message,
		DateTime:  date.UtcDateTimeNow(),
		Trace:     make([]string, 0),
	})
}
