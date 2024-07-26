package domain

import (
	"fmt"
	"os"
	"testing"

	"github.com/danhilltech/100kb.golang/pkg/http"
)

func TestChrome(t *testing.T) {

	chrome, err := startChrome(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer chrome.Shutdown()

	analysis, err := chrome.GetChromeAnalysis("https://caffeine.wiki")

	if err != nil {
		t.Log(err)
		t.Fail()

	}
	analysis.FinalBody = ""
	t.Logf("%+v", analysis)

}

func TestHTTP(t *testing.T) {
	c, err := http.NewClient("/workspaces/100kb.golang/.cache", nil)
	if err != nil {
		t.Fatal(err)
	}
	r, err := c.Get("https://vkc.sh/jeff-geerling-eclipse-experience/")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.StatusCode)
	fmt.Println(r.Header)
}
