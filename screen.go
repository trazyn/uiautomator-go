/**
Screen api related
https://github.com/openatx/uiautomator2#screen-related
*/
package uiautomator

import (
	"encoding/base64"
	"net/http"
)

const (
	ORIENTATION_RIGHT      = "right"
	ORIENTATION_LEFT       = "left"
	ORIENTATION_NATURAL    = "natural"
	ORIENTATION_UPSIDEDOWN = "upsidedown"
)

type (
	ORIENTATION string

	Screenshot struct {
		Type   string
		Base64 string
	}
)

func (ua *UIAutomator) setOrientation(orientation ORIENTATION) error {
	return ua.post(
		&RPCOptions{
			Method: "setOrientation",
			Params: []interface{}{},
		},
		nil,
		nil,
	)
}

/*
Set orientation natural
*/
func (ua *UIAutomator) SetOrientationNatural() error {
	return ua.setOrientation(ORIENTATION_NATURAL)
}

/*
Set orientation upsidedown(not worked)
*/
func (ua *UIAutomator) SetOrientationUpsidedown() error {
	return ua.setOrientation(ORIENTATION_UPSIDEDOWN)
}

/*
Set orientation left
*/
func (ua *UIAutomator) SetOrientationLeft() error {
	return ua.setOrientation(ORIENTATION_LEFT)
}

/*
Set orientation right
*/
func (ua *UIAutomator) SetOrientationRight() error {
	return ua.setOrientation(ORIENTATION_RIGHT)
}

/*
Freeze rotation
*/
func (ua *UIAutomator) FreezeRotation(freeze bool) error {
	return ua.post(
		&RPCOptions{
			Method: "freezeRotation",
			Params: []interface{}{},
		},
		nil,
		nil,
	)
}

/**
Open notification
*/
func (ua *UIAutomator) OpenNotification() error {
	return ua.post(
		&RPCOptions{
			Method: "openNotification",
			Params: []interface{}{},
		},
		nil,
		nil,
	)
}

/**
Open quick settings
*/
func (ua *UIAutomator) OpenQuickSettings() error {
	return ua.post(
		&RPCOptions{
			Method: "openQuickSettings",
			Params: []interface{}{},
		},
		nil,
		nil,
	)
}

/**
Get the UI hierarchy dump content (unicoded).
*/
func (ua *UIAutomator) DumpWindowHierarchy() (string, error) {
	var xml string
	transform := func(payload interface{}, response *http.Response) error {
		xml = payload.(string)
		return nil
	}

	return xml, ua.post(
		&RPCOptions{
			Method: "dumpWindowHierarchy",
			Params: []interface{}{true},
		},
		nil,
		transform,
	)
}

func (ua *UIAutomator) GetScreenshot() (*Screenshot, error) {
	result := &Screenshot{}
	transform := func(data interface{}, response *http.Response) error {
		// Convert to base64
		result.Base64 = base64.StdEncoding.EncodeToString(data.([]byte))
		return nil
	}

	return result, ua.get(
		&RPCOptions{
			URL: "screenshot/0",
		},
		nil,
		transform,
	)
}
