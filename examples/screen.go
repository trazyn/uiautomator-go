package test

import (
	"fmt"
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host: "10.10.20.78",
		Port: 7912,
	})

	xml, err := ua.DumpWindowHierarchy()
	if err != nil {
		panic(err)
	}
	fmt.Println(xml)

	info, err := ua.GetDeviceInfo()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", info)

	size, err := ua.GetWindowSize()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", size)

	app, err := ua.GetCurrentApp()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s/%s\n", app.Package, app.Activity)

	serial, err := ua.GetSerialNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(serial)

	err = ua.WakeUp()
	if err != nil {
		panic(err)
	}

	isWakeup, err := ua.IsWakeUp()
	if err != nil {
		panic(err)
	}
	fmt.Println(isWakeup)

	err = ua.Press("home")
	if err != nil {
		panic(err)
	}

	/*
		err = ua.PressKeyCode(0x07, 0x02)
		if err != nil {
			panic(err)
		}

		// In the home screen click the quick search
		err = ua.DbClick(0.5, 100, 0.1)
		if err != nil {
			panic(err)
		}

		// After open the chrome click the first search tip
		err = ua.Click(0.237, 0.16)
		if err != nil {
			panic(err)
		}

		// Let the search bar move able
		err = ua.LongClick(0.486, 0.096, 0)
		if err != nil {
			panic(err)
		}

		err = ua.Swipe(
			&ug.Position{
				X: 0.836,
				Y: 0.386,
			},

			&ug.Position{
				X: 0.186,
				Y: 0.367,
			},

			0.1,
		)
		if err != nil {
			panic(err)
		}
	*/

	err = ua.GetElementBySelector(
		map[string]interface{}{
			"resourceId": "com.miui.home:id/icon_icon",
			"className":  "android.widget.ImageView",
			"instance":   2,
		},
	)
	if err != nil {
		panic(err)
	}
}
