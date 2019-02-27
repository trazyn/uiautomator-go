package uiautomator

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	VERSION  = "0.0.1"
	BASE_URL = "/jsonrpc/0"

	TIMEOUT                      = 30  // Default timeout(second)
	AUTO_RETRY                   = 5   // Default retry times
	RETRY_DURATION               = 3   // Default retry duration
	WAIT_FOR_EXISTS_MAX_RETRY    = 3   // Default WaitForExistsMaxRetry
	WAIT_FOR_EXISTS_DURATION     = 0.3 // Default WaitForExistsDuration
	WAIT_FOR_DISAPPEAR_MAX_RETRY = 3   // Default WaitForDisappearMaxRetry
	WAIT_FOR_DISAPPEAR_DURATION  = 0.3 // Default WaitForDisappearDuration
)

type (
	RPCOptions struct {
		URL    string
		Method string
		Params []interface{}
	}

	UIAutomator struct {
		config     *Config
		http       *http.Client
		retryTimes int
		size       *WindowSize
	}

	Config struct {
		Host                     string  // Server host
		Port                     int     // Server port
		Timeout                  int     // Timeout(second)
		AutoRetry                int     // Auto retry times, 0 is without retry
		RetryDuration            int     // Retry duration(second)
		WaitForExistsDuration    float32 // Unit second
		WaitForExistsMaxRetry    int     // Max retry times
		WaitForDisappearDuration float32 // Unit second
		WaitForDisappearMaxRetry int     // Max retry times
	}
)

func New(config *Config) *UIAutomator {
	if config == nil {
		panic("New: config can not be null")
	}

	address := net.ParseIP(config.Host)
	if address == nil {
		errMessage := fmt.Sprintf("Incorrect Config.Host: %s", config.Host)
		panic(errMessage)
	} else {
		config.Host = address.String()
	}

	if config.Port <= 0 || config.Port >= 65535 {
		errMessage := fmt.Sprintf("Incorrect Config.Port: %d", config.Port)
		panic(errMessage)
	}

	if config.Timeout < 0 || config.Timeout > 60 {
		config.Timeout = TIMEOUT
	}

	if config.AutoRetry < 0 || config.AutoRetry > 10 {
		config.AutoRetry = AUTO_RETRY
	}

	if config.RetryDuration < 0 || config.RetryDuration > 60 {
		config.RetryDuration = RETRY_DURATION
	}

	if config.WaitForExistsDuration < 0 || config.WaitForExistsDuration > 60 {
		config.WaitForExistsDuration = WAIT_FOR_EXISTS_DURATION
	}

	if config.WaitForExistsMaxRetry < 0 || config.WaitForExistsMaxRetry > 10 {
		config.WaitForExistsMaxRetry = WAIT_FOR_EXISTS_MAX_RETRY
	}

	if config.WaitForDisappearDuration < 0 || config.WaitForDisappearDuration > 60 {
		config.WaitForDisappearDuration = WAIT_FOR_DISAPPEAR_DURATION
	}

	if config.WaitForDisappearMaxRetry < 0 || config.WaitForDisappearMaxRetry > 10 {
		config.WaitForDisappearMaxRetry = WAIT_FOR_DISAPPEAR_MAX_RETRY
	}

	return &UIAutomator{
		config: config,
		http: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		retryTimes: 0,
	}
}

func (ua UIAutomator) GetConfig() *Config {
	return ua.config
}

func (ua *UIAutomator) Ping() (status string, err error) {
	transform := func(response *http.Response) error {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		status = string(responseBody)
		return nil
	}

	err = ua.get(
		&RPCOptions{
			URL: "/ping",
		},
		nil,
		transform,
	)

	if err != nil {
		return
	}

	return
}

func (ua *UIAutomator) caniRetry(err error) bool {
	shouldRetry := true &&
		// Auto retry time should not 0
		ua.config.AutoRetry > 0 &&
		// Retry duration should not 0
		ua.config.RetryDuration > 0 &&
		// Retry time should be less than max auto retry times
		ua.retryTimes < ua.config.AutoRetry

	if shouldRetry {
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return true
			}

		case *url.Error:
			if err.Timeout() {
				return true
			}
		}
	}

	return false
}

func (ua *UIAutomator) execute(request *http.Request, result interface{}, transform interface{}) error {
	for {
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		request.Header.Set("User-Agent", "UIAUTOMATOR/"+VERSION)

		response, err := ua.http.Do(request)
		if err != nil {
			if ua.caniRetry(err) {
				time.Sleep(time.Duration(ua.config.RetryDuration) * time.Second)
				ua.retryTimes++
				continue
			}
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return boom(response)
		}

		// Bypass the body parser
		if transform != nil {
			switch fn := transform.(type) {
			case func(interface{}, *http.Response) error:
				// Pass
			case func(*http.Response) error:
				return fn(response)
			default:
				// Inavlid transform
				transform = nil
			}
		}

		payload, err := parse(response)
		if err != nil {
			return err
		}

		// Everything is ok
		if transform != nil {
			return transform.(func(interface{}, *http.Response) error)(payload, response)
		} else {
			// Decode the JSON result
			if result != nil {
				rawJson, err := json.Marshal(payload)
				if err != nil {
					return err
				}

				if err = json.NewDecoder(bytes.NewBuffer(rawJson)).Decode(&result); err != nil {
					return err
				}
			}
		}
		break
	}

	return nil
}

func (ua *UIAutomator) post(options *RPCOptions, result interface{}, transform interface{}) error {
	requestURL := fmt.Sprintf("http://%s:%d%s", ua.config.Host, ua.config.Port, BASE_URL)
	payload := struct {
		Jsonrpc string        `json:"jsonrpc"`
		ID      string        `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		Jsonrpc: "2.0",
		ID: func() string {
			text := fmt.Sprintf("%s at %u", options.Method, time.Now().Unix())
			hasher := md5.New()
			hasher.Write([]byte(text))
			return hex.EncodeToString(hasher.Sum(nil))
		}(),
		Method: options.Method,
		Params: options.Params,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return ua.execute(request, result, transform)
}

func (ua *UIAutomator) get(options *RPCOptions, result interface{}, transform interface{}) error {
	requestURL := fmt.Sprintf("http://%s:%d/%s", ua.config.Host, ua.config.Port, options.URL)

	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	return ua.execute(request, result, transform)
}

func parse(response *http.Response) (payload interface{}, err error) {
	var RPCReturned struct {
		Error  *UiaError   `json:"error"`
		Result interface{} `json:"result"`
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.Header.Get("Content-Type") != "application/json" {
		// Not an json result use the raw data
		payload = responseBody
		return
	}

	if len(responseBody) == 0 {
		err = fmt.Errorf("%s - empty body", http.StatusText(response.StatusCode))
		return
	}

	err = json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&RPCReturned)
	if err != nil {
		return
	}

	if RPCReturned.Error != nil {
		err = RPCReturned.Error
		return
	}

	payload = RPCReturned.Result
	return
}
