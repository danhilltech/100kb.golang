package domain

import (
	"fmt"
	"os"
	"testing"

	"github.com/andybalholm/cascadia"
	"github.com/danhilltech/100kb.golang/pkg/http"
	"golang.org/x/net/html"
)

func TestChrome(t *testing.T) {

	chrome, err := startChrome(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer chrome.Shutdown()

	analysis, err := chrome.GetChromeAnalysis("https://danhill.is")

	if err != nil {
		t.Log(err)
		t.Fail()

	}
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

func TestBearBlog(t *testing.T) {
	fmt.Println("Getting BearBlog list...")

	c, err := http.NewClient("/workspaces/100kb.golang/.cache", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Get("https://bearblog.dev/discover/?page=0")
	// handle the error if there is one
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	d := 0

	z, err := html.Parse(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	sel, err := cascadia.Parse("span > a")
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range cascadia.QueryAll(z, sel) {

		for _, at := range a.Attr {
			if at.Key == "href" {
				fmt.Println(at.Val)
			}
		}

		// u, err := url.Parse(feed)

		// if err == nil {
		// 	err = engine.Insert(txn, u.Hostname(), feed)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// }
		d++
	}

}
