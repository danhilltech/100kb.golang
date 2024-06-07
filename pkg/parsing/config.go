package parsing

// Everything inside these is gobbeld up into a string
var textTags = []string{"p", "h1", "h2", "h3", "li", "blockquote", "code"}

var badAreas = []string{"nav", "footer", "iframe"}

var internalLinkTags = []string{"a"}

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
