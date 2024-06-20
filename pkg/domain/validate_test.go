package domain

import (
	"fmt"
	"testing"
)

func TestChrome(t *testing.T) {

	chrome, err := startChrome()
	if err != nil {
		t.Fatal(err)
	}
	defer chrome.Shutdown()

	_, err = chrome.GetDomFromChrone("https://danhill.is")
	fmt.Println(err)

	_, err = chrome.GetDomFromChrone("https://dailymail.co.uk")
	fmt.Println(err)
}
