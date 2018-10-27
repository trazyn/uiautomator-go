package uiautomator

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	GatewayError struct {
		Message string
	}
	SessionError struct {
		Message string
	}
	UiaError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func (err *GatewayError) Error() string {
	return err.Message
}

func (err *SessionError) Error() string {
	return err.Message
}

func (err *UiaError) Error() string {
	return err.Message
}

func boom(response *http.Response) error {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusBadGateway {
		return &GatewayError{"Gateway error"}
	}

	if response.StatusCode == http.StatusGone {
		return &SessionError{"App quit or crash"}
	}

	return fmt.Errorf("HTTP Return code is not 200: (%d) [%s]", response.StatusCode, responseBody)
}
