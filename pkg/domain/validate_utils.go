package domain

import (
	_ "embed"
	"net/url"
	"strings"
)

//go:embed data/names.txt
var namesList string

//go:embed data/popular-domains.txt
var popularDomainsList string

func (engine *Engine) identifyURL(baseUrl string) (bool, bool, bool, bool, error) {

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
		if strings.HasSuffix(baseUrlP.Hostname(), domain) {
			popularDomain = true
			break
		}
	}

	urlNews = strings.Contains(strings.ToLower(baseUrlP.Hostname()), "news") || strings.Contains(strings.ToLower(baseUrlP.Hostname()), "daily") || strings.Contains(strings.ToLower(baseUrlP.Hostname()), "standard")

	urlBlog = strings.Contains(strings.ToLower(baseUrlP.Hostname()), "blog") || strings.Contains(strings.ToLower(baseUrlP.Path), "blog")

	return urlHumanName, urlNews, urlBlog, popularDomain, nil
}
