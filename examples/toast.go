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

	ua.Unlock()

	// Show toast
	toast := ua.NewToast()
	toast.Show("Fuckiiiing", 10)
}
