package domain

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestChrome(t *testing.T) {

	chrome, err := startChrome(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer chrome.Shutdown()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		t.Log("starting")
		defer wg.Done()
		body, err := chrome.GetDomFromChrone("https://danhill.is")
		fmt.Println(err)
		t.Log(body)

	}()

	go func() {
		defer wg.Done()
		_, err := chrome.GetDomFromChrone("https://dailymail.co.uk")
		fmt.Println(err)
	}()

	wg.Wait()
}
