/**
https://github.com/openatx/uiautomator2#input-method
*/
package uiautomator

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const _FASTIME = "com.github.uiautomator/.FastInputIME"

var _CODE = map[string]int{
	"go":       2,
	"search":   3,
	"send":     4,
	"next":     5,
	"done":     6,
	"previous": 7,
}

/*
Wait FastInputIME is ready
*/
func (ua *UIAutomator) waitFastinputIME() error {
	r := regexp.MustCompile(`mCurMethodId=([-_./\w]+)`)
	retry := 0

	for {
		if retry > 2 {
			return fmt.Errorf("FastInputIME started failed")
		}

		output, err := ua.Shell([]string{"dumpsys", "input_method"}, 10)
		if err != nil {
			return err
		}

		matchd := r.FindStringSubmatch(output)

		if len(matchd) == 0 || matchd[1] != _FASTIME {
			err := ua.SetFastinputIME(true)
			if err != nil {
				return err
			}

			// Sleep 0.5s
			time.Sleep(time.Duration(500) * time.Millisecond)
			retry++
			continue
		}

		break
	}

	return nil
}

func (ua *UIAutomator) SetFastinputIME(enable bool) error {
	if enable {
		if _, err := ua.Shell([]string{"ime", "enable", _FASTIME}, 5); err != nil {
			return err
		}
		if _, err := ua.Shell([]string{"ime", "set", _FASTIME}, 5); err != nil {
			return err
		}
	} else {
		if _, err := ua.Shell([]string{"ime", "disable", _FASTIME}, 5); err != nil {
			return err
		}
	}
	return nil
}

func (ua *UIAutomator) SendAction(code interface{}) error {
	if err := ua.waitFastinputIME(); err != nil {
		return err
	}

	switch typed := code.(type) {
	case string:
		value, ok := _CODE[typed]
		if !ok {
			return fmt.Errorf("Unknow code: %q", code)
		}
		code = value
	case int:
		// Pass
	default:
		return fmt.Errorf("Unknow code: %q", code)
	}

	if _, err := ua.Shell([]string{"am", "broadcast", "-a", "ADB_EDITOR_CODE", "--ei", "code", strconv.Itoa(code.(int))}, 5); err != nil {
		return err
	}
	return nil
}
