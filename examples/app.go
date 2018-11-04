package main

import (
	"fmt"
	"time"
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host: "10.10.20.66",
		Port: 7912,
	})

	status, err := ua.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("status %v\n", status)

	ua.Unlock()
	ua.AppStart("com.android.chrome")
	time.Sleep(time.Duration(1000) * time.Millisecond)
	ua.AppStop("com.android.chrome")
}
