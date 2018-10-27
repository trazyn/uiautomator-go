package main

import (
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host:      "10.10.20.78",
		Port:      7912,
		AutoRetry: 0,
		Timeout:   10,
	})

	ele, err := ua.GetElementBySelector(
		map[string]interface{}{
			"resourceId": "com.android.chrome:id/url_bar",
		},
	)
	if err != nil {
		panic(err)
	}

	// Set text
	err = ele.SetText("https://www.google.com/")
	if err != nil {
		panic(err)
	}

	// Search
	err = ua.SendAction("search")
	if err != nil {
		panic(err)
	}
}
