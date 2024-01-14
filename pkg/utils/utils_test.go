package utils

import (
	"fmt"
	"testing"
)

func TestURLs(t *testing.T) {
	url := "https://thephd.dev/conformance-should-mean-something-fputc-and-freestanding"

	res, err := URLToPath(url)

	if err != nil {
		t.Fatal(err)
	}

	if res != "thephd.dev/conformance-should-mean-something-fputc-and-freestanding.html" {
		t.Fail()
	}

	url = "https://sonnet.io/posts/hummingbirds/"
	res, err = URLToPath(url)

	if err != nil {
		t.Fatal(err)
	}

	if res != "sonnet.io/posts/hummingbirds.html" {
		fmt.Println(res)
		t.Fail()
	}

	url = "https://sonnet.io/posts/hummingbirds/index.html"
	res, err = URLToPath(url)

	if err != nil {
		t.Fatal(err)
	}

	if res != "sonnet.io/posts/hummingbirds/index.html" {
		fmt.Println(res)
		t.Fail()
	}

}
