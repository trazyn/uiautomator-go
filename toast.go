/**
https://github.com/openatx/uiautomator2#toast
*/
package uiautomator

type Toast struct {
	ua     *UIAutomator
	cached string
}

/*
Get the toast message
TODO: method "getLastToast" not work
*/
func (t *Toast) GetMessage(timeout float32, cachedTime float32, fallback string) (string, error) {
	return fallback, nil
}

/*
Reset the toast cache
TODO: method "getLastToast" not work
*/
func (t *Toast) Reset(message string, duration float32) error {
	return nil
}

/*
Show toast
*/
func (t *Toast) Show(message string, duration float32) error {
	return t.ua.post(
		&RPCOptions{
			Method: "makeToast",
			Params: []interface{}{message, duration * 1000},
		},
		nil,
		nil,
	)
}

/*
Create toast
*/
func (ua *UIAutomator) NewToast() *Toast {
	return &Toast{ua: ua}
}
