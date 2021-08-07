package floodgate

import (
	"context"
	"fmt"
	"net/http"
)

type httpservice struct {
	OriginErr    error  `json:"error"`
	StatusText   string `json:"status"`
	MessageText  string `json:"message"`
	CodeTreshold int    `json:"code_treshold"`
	StatusCode   int    `json:"status_code"`
	Url          string `json:"url"`
	client       Doer
}

// Gets the state error. Nil if there is no error.
func (h *httpservice) Err() error {
	return h.OriginErr
}

// Gets the state status. Returns "healthy" or "error"
func (h *httpservice) Status() string {
	return h.StatusText
}

// Gets the message. Error message if there is error. "ok" if there is no error
func (h *httpservice) Message() string {
	return h.MessageText
}

// Error to be set.
func (h *httpservice) SetError(err error) {
	h.OriginErr = err
}

// Status to be set.
func (h *httpservice) SetStatus(str string) {
	h.StatusText = str
}

// Message to be set.
func (h *httpservice) SetMessage(msg string) {
	h.MessageText = msg
}

// Given context is already with timeout set. No need to set it again.
func (h *httpservice) Check(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, h.Url, nil)

	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= h.CodeTreshold {
		err := fmt.Errorf("service returned status code of '%d'", res.StatusCode)
		return err
	}
	return nil
}

func (h *httpservice) Clone() Checker {
	var c *httpservice
	*c = *h
	return c
}
