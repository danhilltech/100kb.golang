package parsing

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"golang.org/x/net/html"
)

type SimpleNodeType string

type SimpleNode struct {
	Type     SimpleNodeType
	Text     string
	Children []*SimpleNode
}

// Everything inside these is gobbeld up into a string
var textTags = []string{"p", "h1", "h2", "h3", "li", "blockquote"}

func tagIsTextTag(tag string) bool {
	for _, t := range textTags {
		if t == tag {
			return true
		}
	}
	return false
}

var badAreas = []string{"nav", "footer", "iframe", "code", "pre"}

func tagIsGoodArea(tag string) bool {
	for _, t := range badAreas {
		if t == tag {
			return false
		}
	}
	return true
}

var badClassesAndIds = []string{
	"share",
	"widget-area",
	"no-comments",
	"sidebar",
	"sharedaddy",
	"hidden",
	"comments-area",
	"disqus_thread",
	"keep-reading-section",
	"author-box",
	"comment-section",
	"comment",
	"conversation",
	"comment-list",
	"comments",
	"comments-v2",
	"copyright",
	"license",
	"toolbar",
	"twitter-tweet",
	"post-meta",
}

var badWildcards = []string{
	"footer",
	"dropdown",
	"hidden",
}

func tagIsGoodClassOrId(classRaw string) bool {
	class := strings.ToLower(classRaw)

	classes := strings.Split(class, " ")
	for _, t := range badClassesAndIds {

		for _, c := range classes {
			if c == t {
				return false
			}
		}
	}

	for _, t := range badWildcards {
		if strings.Contains(class, t) {
			return false
		}
	}
	return true
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

	}

	if n.Type == html.ElementNode {
		isSafeClass := true
		for _, attr := range n.Attr {
			if attr.Key == "class" || attr.Key == "id" {
				if !tagIsGoodClassOrId(attr.Val) {
					isSafeClass = false
				}
			}
		}

		if !tagIsGoodArea(n.Data) || !isSafeClass {
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
			if tagIsTextTag(p.Data) {
				decendentFromText = true
			}
			if p.Parent == nil {
				break
			}
			p = p.Parent
		}

		if decendentFromText && len(n.Data) > 0 {
			data := []byte(n.Data)
			clean := replaceMultipleWhitespace(data)
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
	if tagIsTextTag(string(node.Type)) {
		workingNode = &SimpleNode{Type: node.Type}
	}

	if node.Type == "text" && workingNode != nil {
		workingNode.Text = workingNode.Text + node.Text
	}

	for _, c := range node.Children {
		walkSimpleNodes(c, workingNode, out)
	}

	if tagIsTextTag(string(node.Type)) {
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

func HtmlToText(htmlBody io.Reader) ([]*serialize.FlatNode, string, string, error) {
	z, err := html.Parse(htmlBody)

	if err != nil {
		return nil, "", "", err
	}

	var title, description string

	rootNode := SimpleNode{Type: "root"}

	walkHtmlNodes(z, &rootNode, 0, &title, &description)

	// Now parse it into strings

	fmt.Println("title", title)

	var simple []*serialize.FlatNode

	walkSimpleNodes(&rootNode, nil, &simple)

	return simple, title, description, nil
}

func (node *SimpleNode) String() string {
	b, _ := json.MarshalIndent(node, "  ", " ")
	return string(b)

}
