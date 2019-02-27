package main

import (
	"fmt"
	ug "uiautomator"
)

func main() {
	ua := ug.New(&ug.Config{
		Host:      "10.10.60.19",
		Port:      7912,
		AutoRetry: 0,
		Timeout:   10,
	})

	eles := make([]*ug.Element, 0)
	ele := ua.GetElementBySelector(map[string]interface{}{"className": "android.widget.ScrollView"})
	fmt.Println(ele.Count())

	ele = ele.Child(map[string]interface{}{"className": "android.view.ViewGroup"})
	count, _ := ele.Count()
	fmt.Println("count:", count)

	height := 139
	for i := 0; i < count; i++ {
		eleItem := ele.Eq(i)
		rect, _ := eleItem.GetRect()
		fmt.Println(rect)
		if rect.Bottom-rect.Top == height {
			eles = append(eles, eleItem)
			fmt.Println(eleItem.GetRect())
			// str := parseElement(eleItem)
			// fmt.Println(str)
			i += 4
		}
	}
	fmt.Println("eles:", len(eles))
	fmt.Println(eles)

	fmt.Println("获取到的元素！！！！！")
	for _, e := range eles {
		rect, _ := e.GetRect()
		fmt.Println("rect:", rect)
	}

	/*
		// Get child element

		ele, err = ele.ChildByText(
			"Clock",
			map[string]interface{}{
				"className": "android.widget.FrameLayout",
			},
		)
	*/

	/*
		// Get element by index

		ele, err = ele.Eq(0)
		if err != nil {
			panic(err)
		}
	*/

	/*
		// Get text

		text, err := ele.GetText()
		if err != nil {
			panic(err)
		}
		fmt.Println(text)
	*/

	/*
		// Set text

		err = ele.SetText("https://www.google.com/")
		if err != nil {
			panic(err)
		}
	*/

	/*
		// Long click

		err = ele.LongClick()
		if err != nil {
			panic(err)
		}
	*/

	/*
		// Swipe element

		err = ele.SwipeLeft()
		if err != nil {
			panic(err)
		}
	*/
}
