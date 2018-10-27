package test

import (
	"fmt"
	ug "uiautomator"
)

func main() {
	client := ug.New(&ug.Config{
		Host: "10.10.20.78",
		Port: 7912,
	})

	xml, err := client.DumpWindowHierarchy()
	if err != nil {
		panic(err)
	}
	fmt.Println(xml)

	info, err := client.GetDeviceInfo()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", info)

	size, err := client.GetWindowSize()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", size)

	app, err := client.GetCurrentApp()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s/%s\n", app.Package, app.Activity)

	serial, err := client.GetSerialNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(serial)

	err = client.WakeUp()
	if err != nil {
		panic(err)
	}

	isWakeup, err := client.IsWakeUp()
	if err != nil {
		panic(err)
	}
	fmt.Println(isWakeup)

	err = client.Press("home")
	if err != nil {
		panic(err)
	}

	/*
		err = client.PressKeyCode(0x07, 0x02)
		if err != nil {
			panic(err)
		}

		// In the home screen click the quick search
		err = client.DbClick(0.5, 100, 0.1)
		if err != nil {
			panic(err)
		}

		// After open the chrome click the first search tip
		err = client.Click(0.237, 0.16)
		if err != nil {
			panic(err)
		}

		// Let the search bar move able
		err = client.LongClick(0.486, 0.096, 0)
		if err != nil {
			panic(err)
		}

		err = client.Swipe(
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

	err = client.GetElementBySelector(
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
