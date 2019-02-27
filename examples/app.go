package main

import (
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host:      "10.10.60.126",
		Port:      7912,
		AutoRetry: 0,
		Timeout:   10,
	})

	ua.Watchman().
		Remove("CIB_RESOLVE_TIMEOUT").
		Register(
			"CIB_RESOLVE_TIMEOUT",
			map[string]interface{}{
				"text": "操作超时，请重新登录",
			},
		).
		Click(
			map[string]interface{}{
				"text": "重新启动",
			},
		)
}
