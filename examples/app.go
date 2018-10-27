package main

import (
	"time"
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host: "10.10.20.78",
		Port: 7912,
	})

	ua.Unlock()
	ua.AppStart("com.android.chrome")
	time.Sleep(time.Duration(1000) * time.Millisecond)
	ua.AppStop("com.android.chrome")
}
