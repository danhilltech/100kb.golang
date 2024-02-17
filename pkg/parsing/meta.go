package parsing

import (
	_ "embed"
	"fmt"
	"net/url"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

//go:embed data/names.txt
var namesList string

//go:embed data/popular-domains.txt
var popularDomainsList string

func (engine *Engine) IdentifyInternalPages(z *html.Node, baseUrl string) (bool, bool, bool, error) {

	var hasAbout, hasBlogRoll, hasWriting bool

	baseUrlP, err := url.Parse(baseUrl)
	if err != nil {
		return false, false, false, fmt.Errorf("could not parse url %s %w", baseUrl, err)
	}

	for _, ic := range internalLinkTags {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return false, false, false, fmt.Errorf("could not create cascadia filter %w", err)
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			for _, attr := range a.Attr {
				if attr.Key == "href" && (strings.HasPrefix(attr.Val, "http") || strings.HasPrefix(attr.Val, "/")) {
					valUrl := strings.TrimSpace(strings.ToLower(attr.Val))

					uP, err := url.Parse(valUrl)
					if err != nil {
						return false, false, false, fmt.Errorf("could not parse url '%s' %w", valUrl, err)
					}
					resolv := baseUrlP.ResolveReference(uP)

					if resolv.Hostname() == baseUrlP.Hostname() && a.FirstChild != nil {
						text := strings.ToLower(a.FirstChild.Data)
						if strings.Contains(text, "about") {
							hasAbout = true
						}
						if strings.Contains(text, "blogroll") {
							hasBlogRoll = true
						}
						if strings.Contains(text, "writing") {
							hasWriting = true
						}
					}

				}
			}

		}
	}

	return hasAbout, hasBlogRoll, hasWriting, nil
}

func (engine *Engine) IdentifyURL(baseUrl string) (bool, bool, bool, bool, error) {

	var urlHumanName, urlNews, urlBlog, popularDomain bool

	baseUrlP, err := url.Parse(baseUrl)
	if err != nil {
		return false, false, false, false, err
	}

	names := strings.Split(namesList, "\n")

	for _, name := range names {
		if strings.Contains(baseUrlP.Hostname(), strings.ToLower(name)) {
			urlHumanName = true
			break
		}
	}

	domains := strings.Split(popularDomainsList, "\n")

	for _, domain := range domains {
		if baseUrlP.Hostname() == strings.ToLower(domain) {
			popularDomain = true
			break
		}
	}

	urlNews = strings.Contains(strings.ToLower(baseUrlP.Hostname()), "news") || strings.Contains(strings.ToLower(baseUrlP.Hostname()), "daily") || strings.Contains(strings.ToLower(baseUrlP.Hostname()), "standard")

	urlBlog = strings.Contains(strings.ToLower(baseUrlP.Hostname()), "blog") || strings.Contains(strings.ToLower(baseUrlP.Path), "blog")

	return urlHumanName, urlNews, urlBlog, popularDomain, nil
}
