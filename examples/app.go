package main

import (
	"fmt"
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host:                  "10.10.60.11",
		Port:                  7912,
		Timeout:               10,
		WaitForExistsMaxRetry: 3,
	})

	status, err := ua.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("status %v\n", status)

	res, err := ua.GetScreenshot()
	if err != nil {
		panic(err)
	}

	fmt.Printf(res.Base64)
}
