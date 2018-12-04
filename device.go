/**
Device api related
https://github.com/openatx/uiautomator2#retrieve-the-device-info
*/
package uiautomator

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type (
	DeviceInfo struct {
		CurrentPackageName string `json:"currentPackageName"`
		DisplayHeight      int    `json:"displayHeight"`
		DisplayWidth       int    `json:"displayWidth"`
		DisplayRotation    int    `json:"displayRotation"`
		DisplaySizeDpX     int    `json:"displaySizeDpX"`
		DisplaySizeDpY     int    `json:"displaySizeDpY"`
		ProductName        string `json:"productName"`
		ScreenOn           bool   `json:"screenOn"`
		SdkInt             int    `json:"sdkInt"`
		NaturalOrientation bool   `json:"naturalOrientation"`
	}

	WindowSize struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	AppInfo struct {
		Activity string `json:"activity"`
		Package  string `json:"package"`
	}
)

/*
Get basic information
*/
func (ua *UIAutomator) GetDeviceInfo() (*DeviceInfo, error) {
	result := &DeviceInfo{}

	return result, ua.post(
		&RPCOptions{
			Method: "deviceInfo",
			Params: []interface{}{},
		},
		result,
		nil,
	)
}

/*
Get window size
*/
func (ua *UIAutomator) GetWindowSize() (*WindowSize, error) {
	var RPCReturned struct {
		Display *WindowSize `json:"display"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	return RPCReturned.Display, ua.get(
		&RPCOptions{
			URL: "info",
		},
		nil,
		transform,
	)
}

/*
Get current app info
*/
func (ua *UIAutomator) GetCurrentApp() (info *AppInfo, err error) {
	output, err := ua.Shell([]string{"dumpsys", "window", "windows"}, 10)
	if err != nil {
		return
	}

	r := regexp.MustCompile(`mCurrentFocus=Window{.*\s+(?P<package>[^\s]+)/(?P<activity>[^\s]+)\}`)
	matched := r.FindStringSubmatch(output)
	res := make(map[string]string)

	for i, name := range r.SubexpNames() {
		if i != 0 && len(name) > 0 {
			res[name] = matched[i]
		}
	}

	info = &AppInfo{
		Package:  res["package"],
		Activity: res["activity"],
	}
	return
}

/*
Get device serial number
*/
func (ua *UIAutomator) GetSerialNumber() (string, error) {
	var RPCReturned struct {
		Serial string `json:"serial"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	return RPCReturned.Serial, ua.get(
		&RPCOptions{
			URL: "info",
		},
		nil,
		transform,
	)
}
