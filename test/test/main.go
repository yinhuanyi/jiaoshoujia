package main

import (
	"fmt"
	"net/url"
)

func main() {
	form := url.Values{}
	form.Add("name", "Robby")
	form.Add("age", "20")
	form.Add("hobby", "running")
	encodeForm := form.Encode()
	fmt.Println(form)
	fmt.Println(encodeForm)
	if decodeForm, err := url.ParseQuery(encodeForm); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(decodeForm)
	}

}
