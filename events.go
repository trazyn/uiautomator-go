/**
Key events api related
https://github.com/openatx/uiautomator2#retrieve-the-device-info
*/

package uiautomator

import (
	"reflect"
	"time"
)

/*
Trun on the screen
*/
func (ua *UIAutomator) WakeUp() error {
	return ua.post(
		&RPCOptions{
			Method: "wakeUp",
			Params: []interface{}{},
		},
		nil,
		nil,
	)
}

/*
Trun off the screen
*/
func (ua *UIAutomator) Sleep() error {
	return ua.post(
		&RPCOptions{
			Method: "sleep",
			Params: []interface{}{},
		},
		nil,
		nil,
	)
}

/*
Check current screen status
*/
func (ua *UIAutomator) checkScreenStatus(wakeUpOeSleep bool) (bool, error) {
	info, err := ua.GetDeviceInfo()
	if err != nil {
		return false, err
	}

	return info.ScreenOn == wakeUpOeSleep, nil
}

/*
Check device is wakeup
*/
func (ua *UIAutomator) IsWakeUp() (res bool, err error) {
	res, err = ua.checkScreenStatus(true)
	return
}

/*
Check device is sleep
*/
func (ua *UIAutomator) IsSleep() (res bool, err error) {
	res, err = ua.checkScreenStatus(false)
	return
}

/*
Press key
*/
func (ua *UIAutomator) Press(key string) error {
	return ua.post(
		&RPCOptions{
			Method: "pressKey",
			Params: []interface{}{key},
		},
		nil,
		nil,
	)
}

/*
Press key code
*/
func (ua *UIAutomator) PressKeyCode(key int, meta interface{}) error {
	params := []interface{}{key}

	if reflect.TypeOf(meta).Kind() == reflect.Int {
		params = append(params, meta)
	}

	return ua.post(
		&RPCOptions{
			Method: "pressKeyCode",
			Params: params,
		},
		nil,
		nil,
	)
}

/*
Unblock the device
*/
func (ua *UIAutomator) Unlock() error {
	var done = make(chan bool)

	// This call will cause blocking, after 1s press home
	go func() {
		time.AfterFunc(
			time.Duration(1000)*time.Millisecond,
			func() {
				done <- true
			},
		)
		ua.Shell(
			[]string{"am start -W -n com.github.uiautomator/.IdentifyActivity -e theme black"},
			0,
		)
	}()

	<-done

	return ua.Press("home")
}
