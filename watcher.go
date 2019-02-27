/**
https://github.com/openatx/uiautomator2#watcher
*/
package uiautomator

type Watcher struct {
	name      string
	ua        *UIAutomator
	selectors []interface{}
}

/*
Create a watcher
*/
func (ua *UIAutomator) Watchman() *Watcher {
	return &Watcher{
		ua:        ua,
		selectors: make([]interface{}, 0),
	}
}

/*
Remove watcher
*/
func (watcher *Watcher) Remove(name string) *Watcher {
	watcher.ua.post(
		&RPCOptions{
			Method: "removeWatcher",
			Params: []interface{}{name},
		},
		nil,
		nil,
	)

	return watcher
}

/*
Add trigger condition
*/
func (watcher *Watcher) Register(name string, selector Selector) *Watcher {
	watcher.name = name
	watcher.selectors = append(watcher.selectors, parseSelector(selector))
	return watcher
}

/*
Listener has triggered and click the target
*/
func (watcher *Watcher) Click(selector Selector) error {
	return watcher.ua.post(
		&RPCOptions{
			Method: "registerClickUiObjectWatcher",
			Params: []interface{}{watcher.name, watcher.selectors, parseSelector(selector)},
		},
		nil,
		nil,
	)
}
