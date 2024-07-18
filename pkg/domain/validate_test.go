package domain

import (
	"os"
	"strings"
	"testing"
)

func TestChrome(t *testing.T) {

	chrome, err := startChrome(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer chrome.Shutdown()

	body, err := chrome.GetDomFromChrone("https://www.theguardian.com/us")

	if len(body) > 100 {
		t.Log("have body")
	}

	if err != nil {
		t.Log(err)
		t.Fail()

	}
	if !strings.Contains(body, "google_ads_") {
		t.Log("no body")
		t.Fail()
	}

}
