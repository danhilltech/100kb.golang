package domain

import "strings"

type ChromeAnalysis struct {
	Requests []*ChromeRequest `json:"requests"`
	TTI      int64            `json:"tti"`
	BeganAt  int64            `json:"-"`
}

type ChromeRequest struct {
	Type string
	Size float64
	URL  string
}

func (c *ChromeAnalysis) LoadsGoogleTagManager() bool {
	for _, req := range c.Requests {
		if req != nil && strings.Contains(req.URL, "googletagmanager.") {
			return true
		}
	}
	return false
}

func (c *ChromeAnalysis) LoadsGoogleAds() bool {
	for _, req := range c.Requests {
		if req != nil && strings.Contains(req.URL, "googlesyndication.") {
			return true
		}
	}
	return false
}

func (c *ChromeAnalysis) LoadsGoogleAdServices() bool {
	for _, req := range c.Requests {
		if req != nil && strings.Contains(req.URL, "googleadservices.") {
			return true
		}
	}
	return false
}

func (c *ChromeAnalysis) LoadsPubmatic() bool {
	for _, req := range c.Requests {
		if req != nil && strings.Contains(req.URL, "pubmatic.") {
			return true
		}
	}
	return false
}

func (c *ChromeAnalysis) LoadsTwitterAds() bool {
	for _, req := range c.Requests {
		if req != nil && strings.Contains(req.URL, "ads-twitter.") {
			return true
		}
	}
	return false
}

func (c *ChromeAnalysis) LoadsAmazonAds() bool {
	for _, req := range c.Requests {
		if req != nil && strings.Contains(req.URL, "amazon-adsystem.") {
			return true
		}
	}
	return false
}

func (c *ChromeAnalysis) TotalNetworkRequests() int {
	var out int
	for _, req := range c.Requests {
		if req != nil {
			out++
		}
	}
	return out
}

func (c *ChromeAnalysis) TotalScriptRequests() int {
	var out int
	for _, req := range c.Requests {
		if req != nil && req.Type == "Script" {
			out++
		}
	}
	return out
}

func (c *ChromeAnalysis) TotalCSSRequests() int {
	var out int
	for _, req := range c.Requests {
		if req != nil && req.Type == "Stylesheet" {
			out++
		}
	}
	return out
}

func (c *ChromeAnalysis) TotalWeight() int {
	var out int
	for _, req := range c.Requests {
		if req != nil {
			out += int(req.Size)
		}
	}
	return out
}

func (c *ChromeAnalysis) TotalScriptWeight() int {
	var out int
	for _, req := range c.Requests {
		if req != nil && req.Type == "Script" {
			out += int(req.Size)
		}
	}
	return out
}

func (c *ChromeAnalysis) TotalCSSWeight() int {
	var out int
	for _, req := range c.Requests {
		if req != nil && req.Type == "Stylesheet" {
			out += int(req.Size)
		}
	}
	return out
}

func (c *ChromeAnalysis) TotalDocumentWeight() int {
	var out int
	for _, req := range c.Requests {
		if req != nil && req.Type == "Document" {
			out += int(req.Size)
		}
	}
	return out
}
