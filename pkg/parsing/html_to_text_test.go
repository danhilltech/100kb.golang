package parsing

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func TestBasic(t *testing.T) {

	var pricingHtml string = `
<div class="card mb-4 box-shadow">
	<div class="card-header">
		<h4 class="my-0 font-weight-normal">Free</h4>
	</div>
	<div class="card-body">
		<h1 class="card-title pricing-card-title">$0/mo</h1>
		<ul class="list-unstyled mt-3 mb-4">
			<li>10 users included</li>
			<li>2 GB of storage</li>
			<li><a href="https://example.com">See more</a></li>
		</ul>
	</div>
</div>

<div class="card mb-4 box-shadow">
	<div class="card-header">
		<h4 class="my-0 font-weight-normal">Pro</h4>
	</div>
	<div class="card-body">
		<h1 class="card-title pricing-card-title">$15/mo</h1>
		<ul class="list-unstyled mt-3 mb-4">
			<li>20 users included</li>
			<li>10 GB of storage</li>
			<li><a href="https://example.com">See more</a></li>
		</ul>
	</div>
</div>

<div class="card mb-4 box-shadow">
	<div class="card-header">
		<h4 class="my-0 font-weight-normal">Enterprise</h4>
	</div>
	<div class="card-body">
		<h1 class="card-title pricing-card-title">$29/mo</h1>
		<ul class="list-unstyled mt-3 mb-4">
			<li>30 users included</li>
			<li>15 GB of storage</li>
			<li><a>See more</a></li>
		</ul>
	</div>
</div>
`

	doc, err := html.Parse(strings.NewReader(pricingHtml))
	if err != nil {
		t.Fatal(err)
	}

	sel, err := cascadia.Parse(".card-title")
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range cascadia.QueryAll(doc, sel) {
		t.Log(a.Data)
		a.Attr = append(a.Attr, html.Attribute{Key: "ignore", Val: "yes"})
	}
	var out bytes.Buffer
	w := io.Writer(&out)
	html.Render(w, doc)

	t.Log(out.String())
}
