/**
Selector is a handy mechanism to identify a specific UI object in the current window.
https://github.com/openatx/uiautomator2#selector
*/
package uiautomator

import (
	"encoding/json"
	"net/http"
	"time"
)

type (
	Selector map[string]interface{}

	Element struct {
		ua       *UIAutomator
		position *Position
		selector Selector
	}

	ElementRect struct {
		Bottom int `json:"bottom"`
		Left   int `json:"left"`
		Right  int `json:"right"`
		Top    int `json:"top"`
	}

	ElementInfo struct {
		ContentDescription string       `json:"contentDescription"`
		Checked            bool         `json:"checked"`
		Scrollable         bool         `json:"scrollable"`
		Text               string       `json:"text"`
		PackageName        string       `json:"packageName"`
		Selected           bool         `json:"selected"`
		Enabled            bool         `json:"enabled"`
		ClassName          string       `json:"className"`
		Focused            bool         `json:"focused"`
		Focusable          bool         `json:"focusable"`
		Clickable          bool         `json:"clickable"`
		ChileCount         int          `json:"chileCount"`
		LongClickable      bool         `json:"longClickable"`
		Checkable          bool         `json:"checkable"`
		Bounds             *ElementRect `json:"bounds"`
		VisibleBounds      *ElementRect `json:"visibleBounds"`
	}
)

var _MASK = map[string]int{
	"text":                  0x01,       // MASK_TEXT,
	"textContains":          0x02,       // MASK_TEXTCONTAINS,
	"textMatches":           0x04,       // MASK_TEXTMATCHES,
	"textStartsWith":        0x08,       // MASK_TEXTSTARTSWITH,
	"className":             0x10,       // MASK_CLASSNAME
	"classNameMatches":      0x20,       // MASK_CLASSNAMEMATCHES
	"description":           0x40,       // MASK_DESCRIPTION
	"descriptionContains":   0x80,       // MASK_DESCRIPTIONCONTAINS
	"descriptionMatches":    0x0100,     // MASK_DESCRIPTIONMATCHES
	"descriptionStartsWith": 0x0200,     // MASK_DESCRIPTIONSTARTSWITH
	"checkable":             0x0400,     // MASK_CHECKABLE
	"checked":               0x0800,     // MASK_CHECKED
	"clickable":             0x1000,     // MASK_CLICKABLE
	"longClickable":         0x2000,     // MASK_LONGCLICKABLE,
	"scrollable":            0x4000,     // MASK_SCROLLABLE,
	"enabled":               0x8000,     // MASK_ENABLED,
	"focusable":             0x010000,   // MASK_FOCUSABLE,
	"focused":               0x020000,   // MASK_FOCUSED,
	"selected":              0x040000,   // MASK_SELECTED,
	"packageName":           0x080000,   // MASK_PACKAGENAME,
	"packageNameMatches":    0x100000,   // MASK_PACKAGENAMEMATCHES,
	"resourceId":            0x200000,   // MASK_RESOURCEID,
	"resourceIdMatches":     0x400000,   // MASK_RESOURCEIDMATCHES,
	"index":                 0x800000,   // MASK_INDEX,
	"instance":              0x01000000, // MASK_INSTANCE,
}

/*
Get element info
*/
func (ele Element) GetInfo() (*ElementInfo, error) {
	var RPCReturned ElementInfo

	if err := ele.ua.post(
		&RPCOptions{
			Method: "objInfo",
			Params: []interface{}{getParams(ele.selector)},
		},
		&RPCReturned,
		nil,
	); err != nil {
		return nil, err
	}

	return &RPCReturned, nil
}

/*
Get Widget rect bounds
*/
func (ele Element) GetRect() (rect *ElementRect, err error) {
	info, err := ele.GetInfo()
	if err != nil {
		return
	}

	rect = info.Bounds
	if rect == nil {
		rect = info.VisibleBounds
	}
	return
}

/*
Get Widget center point
*/
func (ele Element) Center(offset *Position) (*Position, error) {
	rect, err := ele.GetRect()
	if err != nil {
		return nil, err
	}

	lx, ly, rx, ry := rect.Left, rect.Top, rect.Right, rect.Bottom
	width, height := rx-lx, ry-ly

	if offset == nil {
		offset = &Position{0.5, 0.5}
	}

	abs := &Position{}
	abs.X = float32(lx) + float32(width)*offset.X
	abs.Y = float32(ly) + float32(height)*offset.Y
	return abs, nil
}

/*
Get the count
*/
func (ele Element) Count() (int, error) {
	var RPCReturned struct {
		Result int `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	return RPCReturned.Result, ele.ua.post(
		&RPCOptions{
			Method: "count",
			Params: []interface{}{getParams(ele.selector)},
		},
		nil,
		transform,
	)
}

/*
Clone the element
*/
func (ele Element) Clone() *Element {
	copied := ele

	copied.selector = make(Selector)

	for k, v := range ele.selector {
		copied.selector[k] = v

		if k == "childOrSiblingSelector" {
			nested := make([]interface{}, 0)
			for _, selector := range v.([]interface{}) {
				nested = append(nested, parseSelector(selector.(Selector)))
			}
			copied.selector[k] = nested
		}
	}

	return &copied
}

/*
Get the instance via index
*/
func (ele Element) Eq(index int) *Element {
	copied := ele.Clone()

	// Check is a child selector
	childOrSiblingSelector := copied.selector["childOrSiblingSelector"].([]interface{})
	lastSelectorIndex := len(childOrSiblingSelector)

	if lastSelectorIndex > 0 {
		// Get the child selector
		lastSelector := childOrSiblingSelector[lastSelectorIndex-1].(Selector)

		// Update the child selector
		lastSelector["instance"] = index
		newSelector := parseSelector(lastSelector)
		childOrSiblingSelector[lastSelectorIndex-1] = newSelector
	} else {
		// Update the selector
		copied.selector["instance"] = index

		// Rebuild the selector
		newSelector := parseSelector(copied.selector)
		copied.selector = newSelector
	}

	return copied
}

/*
Check if the specific UI object exists
*/
func (ele Element) WaitForExists(duration float32, maxRetry int) error {
	if duration < 0 || duration > 60 {
		duration = WAIT_FOR_EXISTS_DURATION
	}

	if maxRetry < 0 || maxRetry > 10 {
		maxRetry = WAIT_FOR_EXISTS_MAX_RETRY
	}

	return ele.wait(duration, maxRetry, true)
}

/*
Wait the specific UI object disappear
*/
func (ele Element) WaitUntilGone(duration float32, maxRetry int) error {
	if duration < 0 || duration > 60 {
		duration = WAIT_FOR_DISAPPEAR_DURATION
	}

	if maxRetry < 0 || maxRetry > 10 {
		maxRetry = WAIT_FOR_DISAPPEAR_MAX_RETRY
	}

	return ele.wait(duration, maxRetry, false)
}

/*
Wait element exists or gone
*/
func (ele Element) wait(duration float32, maxRetry int, exists bool) error {
	var (
		err    error
		retry  int
		method string
	)

	config := ele.ua.GetConfig()

	if exists {
		method = "waitForExists"
	} else {
		method = "waitUntilGone"
	}

	for {
		var RPCReturned struct {
			Result bool `json:"result"`
		}
		transform := func(response *http.Response) error {
			err := json.NewDecoder(response.Body).Decode(&RPCReturned)
			if err != nil {
				return err
			}
			return nil
		}

		err = ele.ua.post(
			&RPCOptions{
				Method: method,
				Params: []interface{}{getParams(ele.selector), config.Timeout * 1000},
			},
			nil,
			transform,
		)

		if err != nil || RPCReturned.Result == false {
			retry++

			if retry < maxRetry {
				time.Sleep(time.Duration(duration*1000) * time.Millisecond)
				continue
			}

			err = &UiaError{
				Code:    -32002,
				Message: "Element not found",
			}

			break
		}

		// It's ok
		break
	}

	return err
}

/*
Swipe the element
*/
func (ele Element) swipe(direction string) error {
	config := ele.ua.GetConfig()
	if err := ele.WaitForExists(config.WaitForExistsDuration, config.WaitForExistsMaxRetry); err != nil {
		return err
	}
	rect, err := ele.GetRect()
	if err != nil {
		return err
	}

	lx, ly, rx, ry := rect.Left, rect.Top, rect.Right, rect.Bottom
	cx, cy := (lx+rx)/2, (ly+ry)/2

	switch direction {
	case "up":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(cx), Y: float32(ly)},
			20,
		)
	case "down":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(cx), Y: float32(ry - 1)},
			20,
		)
	case "left":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(lx), Y: float32(cy)},
			20,
		)
	case "right":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(rx - 1), Y: float32(cy)},
			20,
		)
	}

	return nil
}

/*
Swipe to up
*/
func (ele *Element) SwipeUp() error {
	return ele.swipe("up")
}

/*
Swipe to down
*/
func (ele *Element) SwipeDown() error {
	return ele.swipe("down")
}

/*
Swipe to left
*/
func (ele *Element) SwipeLeft() error {
	return ele.swipe("left")
}

/*
Swipe to right
*/
func (ele *Element) SwipeRight() error {
	return ele.swipe("right")
}

/*
Click on the screen
*/
func (ele *Element) Click(offset *Position) error {
	config := ele.ua.GetConfig()
	if err := ele.WaitForExists(config.WaitForExistsDuration, config.WaitForExistsMaxRetry); err != nil {
		return err
	}

	return ele.ClickNoWait(offset)
}

func (ele *Element) ClickNoWait(offset *Position) error {
	abs, err := ele.Center(offset)
	if err != nil {
		return err
	}

	return ele.ua.Click(abs)
}

/*
Screen scroll up
*/
func (ele *Element) ScrollUp(step int) error {
	if err := ele.ua.post(
		&RPCOptions{
			Method: "scrollForward",
			Params: []interface{}{ele.selector, true, step},
		},
		nil,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
Screen scroll down
*/
func (ele *Element) ScrollDown(step int) error {
	if err := ele.ua.post(
		&RPCOptions{
			Method: "scrollBackward",
			Params: []interface{}{ele.selector, true, step},
		},
		nil,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
Screen scroll to beginning
*/
func (ele *Element) ScrollToBeginning() error {
	if err := ele.ua.post(
		&RPCOptions{
			Method: "flingBackward",
			Params: []interface{}{ele.selector, true},
		},
		nil,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
Screen scroll to end
*/
func (ele *Element) ScrollToEnd() error {
	if err := ele.ua.post(
		&RPCOptions{
			Method: "scrollToEnd",
			Params: []interface{}{ele.selector, true, 500, 20},
		},
		nil,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
Screen scroll to selector
*/
func (ele *Element) ScrollTo(selector Selector) error {
	selector = parseSelector(selector)

	if err := ele.ua.post(
		&RPCOptions{
			Method: "scrollTo",
			Params: []interface{}{ele.selector, selector, true},
		},
		nil,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
Long click on the element
*/
func (ele *Element) LongClick() error {
	config := ele.ua.GetConfig()
	if err := ele.WaitForExists(config.WaitForExistsDuration, config.WaitForExistsMaxRetry); err != nil {
		return err
	}

	abs, err := ele.Center(nil)
	if err != nil {
		return err
	}

	return ele.ua.LongClick(abs, 0)
}

/*
Get the children or grandchildren
*/
func (ele Element) Child(selector Selector) *Element {
	copied := ele.Clone()

	selector = parseSelector(selector)

	var (
		childOrSibling         = copied.selector["childOrSibling"]
		childOrSiblingSelector = copied.selector["childOrSiblingSelector"]
	)

	if childOrSibling == nil {
		childOrSibling = make([]interface{}, 1)
		childOrSibling = append(childOrSibling.([]interface{}), "child")
	}

	if childOrSiblingSelector == nil {
		childOrSiblingSelector = make([]interface{}, 1)
		childOrSiblingSelector = append(childOrSiblingSelector.([]interface{}), selector)
	}

	childOrSibling = append(childOrSibling.([]interface{}), "child")
	childOrSiblingSelector = append(childOrSiblingSelector.([]interface{}), selector)

	copied.selector["childOrSibling"] = childOrSibling
	copied.selector["childOrSiblingSelector"] = childOrSiblingSelector
	return copied
}

func (ele *Element) childByMethod(keywords string, method string, selector Selector) (*Element, error) {
	var RPCReturned struct {
		Result string `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	selector = parseSelector(selector)

	if err := ele.ua.post(
		&RPCOptions{
			Method: method,
			Params: []interface{}{ele.selector, selector, keywords, true},
		},
		nil,
		transform,
	); err != nil {
		return nil, err
	}

	ele.selector = map[string]interface{}{"__UID": RPCReturned.Result}
	return ele, nil
}

func (ele *Element) ChildByText(keywords string, selector Selector) (*Element, error) {
	return ele.childByMethod(keywords, "childByText", selector)
}

func (ele *Element) ChildByDescription(keywords string, selector Selector) (*Element, error) {
	return ele.childByMethod(keywords, "childByDescription", selector)
}

/*
Get the sibling
*/
func (ele *Element) Sibling(selector Selector) (*Element, error) {
	selector = parseSelector(selector)

	ele.selector["childOrSibling"] = []interface{}{"sibling"}
	ele.selector["childOrSiblingSelector"] = []interface{}{selector}

	return ele, nil
}

/*
Get widget text
*/
func (ele Element) GetText() (string, error) {
	config := ele.ua.GetConfig()
	if err := ele.WaitForExists(config.WaitForExistsDuration, config.WaitForExistsMaxRetry); err != nil {
		return "", err
	}

	return ele.GetTextNoWait()
}

func (ele Element) GetTextNoWait() (string, error) {
	var RPCReturned struct {
		Result string `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	return RPCReturned.Result, ele.ua.post(
		&RPCOptions{
			Method: "getText",
			Params: []interface{}{getParams(ele.selector)},
		},
		nil,
		transform,
	)
}

/*
Set widget text
*/
func (ele Element) SetText(text string) error {
	config := ele.ua.GetConfig()
	if err := ele.WaitForExists(config.WaitForExistsDuration, config.WaitForExistsMaxRetry); err != nil {
		return err
	}

	return ele.ua.post(
		&RPCOptions{
			Method: "setText",
			Params: []interface{}{getParams(ele.selector), text},
		},
		nil,
		nil,
	)
}

/*
Clear the widget text
*/
func (ele Element) ClearText() error {
	config := ele.ua.GetConfig()
	if err := ele.WaitForExists(config.WaitForExistsDuration, config.WaitForExistsMaxRetry); err != nil {
		return err
	}

	return ele.ua.post(
		&RPCOptions{
			Method: "clearTextField",
			Params: []interface{}{getParams(ele.selector)},
		},
		nil,
		nil,
	)
}

/*
Query the UI element by selector
*/
func (ua *UIAutomator) GetElementBySelector(selector Selector) (ele *Element) {
	ele = &Element{ua: ua}

	selector = parseSelector(selector)

	ele.selector = selector
	return
}

func parseSelector(selector Selector) Selector {
	res := make(Selector)

	// Params initalization
	res["mask"] = selector["mask"]
	res["childOrSibling"] = selector["childOrSibling"]
	res["childOrSiblingSelector"] = selector["childOrSiblingSelector"]

	if res["mask"] == nil {
		res["mask"] = 0
	}

	if res["childOrSibling"] == nil {
		res["childOrSibling"] = []interface{}{}
	}

	if res["childOrSiblingSelector"] == nil {
		res["childOrSiblingSelector"] = []interface{}{}
	}

	for k, v := range selector {
		if selectorMask, ok := _MASK[k]; ok {
			res[k] = v
			res["mask"] = res["mask"].(int) | selectorMask
		}
	}

	return res
}

func getParams(selector Selector) interface{} {
	if uid, ok := selector["__UID"]; ok {
		return uid
	}
	return selector
}
