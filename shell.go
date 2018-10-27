package uiautomator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (ua *UIAutomator) Shell(command []string, timeout int) (output string, err error) {
	requestURL := fmt.Sprintf("http://%s:%d/shell", ua.config.Host, ua.config.Port)
	response, err := http.PostForm(
		requestURL,
		url.Values{
			"command": {strings.Join(command, " ")},
			"timeout": {strconv.Itoa(timeout)},
		},
	)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = boom(response)
		return
	}

	var ShellReturned struct {
		ExitCode int    `json:"exitCode"`
		Output   string `json:"output"`
	}

	err = json.NewDecoder(response.Body).Decode(&ShellReturned)
	if err != nil {
		return
	}
	if ShellReturned.ExitCode != 0 {
		err = &UiaError{
			Code:    ShellReturned.ExitCode,
			Message: fmt.Sprint("Failed to execute command: %s", command),
		}

		return
	}

	output = ShellReturned.Output
	return
}
