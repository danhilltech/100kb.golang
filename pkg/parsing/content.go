package parsing

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"golang.org/x/net/html"
)

type SimpleNodeType string

type SimpleNode struct {
	Type     SimpleNodeType
	Text     string
	Children []*SimpleNode
}

type ParseAnalysis struct {
	Ids     []string
	Classes []string
	Urls    []string
	Links   []string

	BadUrls       []string
	BadElements   []string
	BadLinkTitles []string
}

var whitespaceTable = [256]bool{
	// ASCII
	false, false, false, false, false, false, false, false,
	false, true, true, false, true, true, false, false, // tab, new line, form feed, carriage return
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	true, false, false, false, false, false, false, false, // space
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	// non-ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}

// IsWhitespace returns true for space, \n, \r, \t, \f.
func isWhitespace(c byte) bool {
	return whitespaceTable[c]
}

var newlineTable = [256]bool{
	// ASCII
	false, false, false, false, false, false, false, false,
	false, false, true, false, false, true, false, false, // new line, carriage return
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	// non-ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}

// IsNewline returns true for \n, \r.
func isNewline(c byte) bool {
	return newlineTable[c]
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func isIncludeNode(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "data-action" && attr.Val == "include" {
			return true
		}
	}
	if n.Parent != nil {
		return isIncludeNode(n.Parent)
	}

	return false
}

func tagIsTextTag(t SimpleNodeType) bool {
	for _, c := range textTags {
		if c == string(t) {
			return true
		}
	}
	return false
}

func extractMetaProperty(t *html.Node, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "name" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}

func replaceMultipleWhitespace(b []byte) []byte {
	j, k := 0, 0 // j is write position, k is start of next text section
	for i := 0; i < len(b); i++ {
		if isWhitespace(b[i]) {
			start := i
			newline := isNewline(b[i])
			i++
			for ; i < len(b) && isWhitespace(b[i]); i++ {
				if isNewline(b[i]) {
					newline = true
				}
			}
			if newline {
				b[start] = ' '
			} else {
				b[start] = ' '
			}
			if 1 < i-start { // more than one whitespace
				if j == 0 {
					j = start + 1
				} else {
					j += copy(b[j:], b[k:start+1])
				}
				k = i
			}
		}
	}
	if j == 0 {
		return b
	} else if j == 1 { // only if starts with whitespace
		b[k-1] = b[0]
		return b[k-1:]
	} else if k < len(b) {
		j += copy(b[j:], b[k:])
	}
	return b[:j]
}

func walkHtmlNodes(n *html.Node, b *SimpleNode, depth int, title *string, description *string) {
	if isTitleElement(n) {
		if n.FirstChild != nil {
			*title = n.FirstChild.Data
		}
	}

	if n.Data == "meta" {
		desc, ok := extractMetaProperty(n, "description")
		if ok && *description == "" {
			*description = desc
		}

		descOg, ok := extractMetaProperty(n, "og:description")
		if ok && *description == "" {
			*description = descOg
		}

		titleOg, ok := extractMetaProperty(n, "og:title")
		if ok {
			*title = titleOg
		}

	}

	if n.Type == html.ElementNode {
		isSafeClass := true
		for _, attr := range n.Attr {
			if attr.Key == "data-action" {
				if attr.Val == "skip" || attr.Val == "block" {
					isSafeClass = false
				}
			}
		}

		if !isSafeClass {
			return
		}
		// Create the new node
		nB := SimpleNode{
			Type: SimpleNodeType(n.Data),
		}

		b.Children = append(b.Children, &nB)

		b = &nB

	}
	if n.Type == html.TextNode {
		decendentFromText := false

		p := n.Parent
		for {
			if isIncludeNode(p) {
				decendentFromText = true
			}
			if p.Parent == nil {
				break
			}
			p = p.Parent
		}

		if decendentFromText && len(n.Data) > 0 {
			data := []byte(n.Data)
			clean := bytes.TrimSpace(replaceMultipleWhitespace(data))
			newNode := SimpleNode{
				Text: string(clean),
				Type: "text",
			}

			b.Children = append(b.Children, &newNode)
		}
	}
	nextDepth := depth + 1

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkHtmlNodes(c, b, nextDepth, title, description)
	}
}

func walkSimpleNodes(node *SimpleNode, workingNode *SimpleNode, out *[]*serialize.FlatNode) {
	if tagIsTextTag(node.Type) {
		workingNode = &SimpleNode{Type: node.Type}
	}

	if node.Type == "text" && workingNode != nil {
		workingNode.Text = workingNode.Text + node.Text
	}

	for _, c := range node.Children {
		walkSimpleNodes(c, workingNode, out)
	}

	if tagIsTextTag(node.Type) {
		txt := strings.TrimSpace(workingNode.Text)
		if len(txt) > 5 {
			txt := replaceMultipleWhitespace([]byte(txt))

			flat := &serialize.FlatNode{
				Type: string(workingNode.Type),
				Text: string(txt),
			}
			*out = append(*out, flat)
		}
	}
}

func (engine *Engine) HtmlToText(z *html.Node) ([]*serialize.FlatNode, string, string, error) {

	var title, description string

	rootNode := SimpleNode{Type: "root"}

	walkHtmlNodes(z, &rootNode, 0, &title, &description)

	var simple []*serialize.FlatNode

	walkSimpleNodes(&rootNode, nil, &simple)

	return simple, title, description, nil
}

func (node *SimpleNode) String() string {
	b, _ := json.MarshalIndent(node, "  ", " ")
	return string(b)

}

func (engine *Engine) IdentifyElements(z *html.Node, baseUrl string, adblock *AdblockEngine) (*ParseAnalysis, error) {

	parseAnalysis := ParseAnalysis{
		Ids:           make([]string, 0),
		Classes:       make([]string, 0),
		Urls:          make([]string, 0),
		Links:         make([]string, 0),
		BadUrls:       make([]string, 0),
		BadElements:   make([]string, 0),
		BadLinkTitles: make([]string, 0),
	}

	walkHtmlNodesAndIdentify(z, &parseAnalysis)

	badIdsAndClasses, badUrls, err := adblock.Filter(parseAnalysis.Ids, parseAnalysis.Classes, parseAnalysis.Urls, baseUrl)
	if err != nil {
		return nil, err
	}

	parseAnalysis.BadUrls = badUrls
	parseAnalysis.BadElements = badIdsAndClasses

	for _, ic := range badIdsAndClasses {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return nil, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "block"})
		}
	}

	for _, ic := range badClassesAndIds {
		sel, err := cascadia.Parse("#" + ic)
		if err != nil {
			return nil, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "skip"})
		}
	}

	for _, ic := range badClassesAndIds {
		sel, err := cascadia.Parse("." + ic)
		if err != nil {
			return nil, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "skip"})
		}
	}

	for _, ic := range badAreas {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return nil, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "skip"})
		}
	}

	for _, ic := range textTags {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return nil, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			hasLabel := false
			for _, ae := range a.Attr {
				if ae.Key == "data-action" {
					hasLabel = true
				}
			}
			if !hasLabel {
				a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "include"})
			}
		}
	}

	return &parseAnalysis, nil
}

func walkHtmlNodesAndIdentify(n *html.Node, pa *ParseAnalysis) {

	if n.Type == html.ElementNode {

		if n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					pa.Links = append(pa.Links, attr.Val)
				}
			}
			if n.FirstChild != nil {
				title := strings.ToLower(n.FirstChild.Data)

				for _, b := range badLinkTitles {
					if strings.Contains(title, b) {
						pa.BadLinkTitles = append(pa.BadLinkTitles, title)
					}
				}

			}
		}

		for _, attr := range n.Attr {
			if attr.Key == "class" {
				nc := strings.Split(attr.Val, " ")
				pa.Classes = append(pa.Classes, nc...)
			}
			if attr.Key == "id" {
				pa.Ids = append(pa.Ids, attr.Val)
			}
			if attr.Key == "href" {
				pa.Urls = append(pa.Urls, attr.Val)
			}
			if attr.Key == "src" {
				pa.Urls = append(pa.Urls, attr.Val)
			}

		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkHtmlNodesAndIdentify(c, pa)
	}
}
