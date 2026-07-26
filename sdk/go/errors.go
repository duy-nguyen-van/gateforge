package gateforge

import (
	"fmt"
	"net/http"
)

// Meta matches the GateForge IAM JSON envelope metadata
// ({"meta": {...}, "data": ...}).
type Meta struct {
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
	Code      int    `json:"code,omitempty"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	Total     int64  `json:"total,omitempty"`
}

// Envelope is the standard GateForge IAM response wrapper.
type Envelope[T any] struct {
	Meta Meta `json:"meta"`
	Data T    `json:"data"`
}

// APIError is returned when an HTTP response indicates failure.
// It surfaces Meta.error_code / Meta.message from the envelope when present.
type APIError struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Meta       Meta
	Response   *http.Response
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("gateforge: HTTP %d %s: %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("gateforge: HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("gateforge: HTTP %d", e.StatusCode)
}

// NewAPIError builds an APIError from status and Meta fields.
func NewAPIError(statusCode int, meta Meta, resp *http.Response) *APIError {
	msg := meta.Message
	if msg == "" {
		msg = http.StatusText(statusCode)
	}
	return &APIError{
		StatusCode: statusCode,
		ErrorCode:  meta.ErrorCode,
		Message:    msg,
		Meta:       meta,
		Response:   resp,
	}
}
