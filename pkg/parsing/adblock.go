package parsing

/*
#cgo LDFLAGS: -L../../lib -lgoadblock
#include "../../lib/goadblock.h"
*/
import "C"
import (
	"embed"
	_ "embed"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"unsafe"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/proto"
)

type AdblockEngine struct {
	engine unsafe.Pointer
	mutex  sync.Mutex
}

//go:embed data/adblock/*.txt
var adbLists embed.FS

func NewAdblockEngine() (*AdblockEngine, error) {

	req := RuleGroups{}

	files, err := adbLists.ReadDir("data/adblock")
	if err != nil {
		return nil, err
	}

	var cnt int

	for _, list := range files {
		contents, err := adbLists.ReadFile(fmt.Sprintf("data/adblock/%s", list.Name()))
		if err != nil {
			return nil, err
		}
		rules := strings.Split(string(contents), "\n")

		cnt += len(rules)

		rs := Rules{}
		rs.Rules = rules

		req.Filters = append(req.Filters, &rs)
	}

	fmt.Printf("Adblock loading %d rules...\n", cnt)
	defer fmt.Printf("Adblock loaded\n")

	reqBytes, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	reqSize := uintptr(len(reqBytes))

	creqSize := unsafe.Pointer(&reqSize)
	reqPtr := unsafe.Pointer(&reqBytes[0])

	a := C.new_adblock((*C.uchar)(reqPtr), (*C.size_t)(creqSize))

	return &AdblockEngine{engine: unsafe.Pointer(a)}, nil

}

func (engine *AdblockEngine) Close() {
	C.drop_adblock((*C.AdblockEngine)(engine.engine))
}

func (engine *AdblockEngine) Filter(ids []string, classes []string, urls []string, baseUrl string) ([]string, []string, error) {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()

	req := FilterRequest{}

	baseUrlP, err := url.Parse(baseUrl)
	if err != nil {
		return nil, nil, err
	}
	if baseUrlP.Scheme == "" || baseUrlP.Hostname() == "" {
		return nil, nil, fmt.Errorf("no scheme or hostname found")
	}

	urlsClean := []string{}

	for _, uRaw := range urls {

		u := strings.TrimSpace(uRaw)
		if strings.HasPrefix(u, "mailto:") {
			continue
		}
		if strings.HasPrefix(u, "about:blank") {
			continue
		}
		if strings.HasPrefix(u, "data:") {
			continue
		}
		if strings.HasPrefix(u, "javascript:") {
			continue
		}
		if strings.HasPrefix(u, "tel:") {
			continue
		}
		if strings.HasPrefix(u, "file:") {
			continue
		}
		if strings.HasPrefix(u, "sms:") {
			continue
		}

		if strings.HasPrefix(u, "feed:") {
			continue
		}

		if strings.HasPrefix(u, "skype:") {
			continue
		}

		uP, err := url.Parse(u)
		if err != nil {
			continue
		}
		resolv := baseUrlP.ResolveReference(uP)
		urlsClean = append(urlsClean, resolv.String())
	}

	req.Classes = classes
	req.Ids = ids
	req.Urls = urlsClean
	req.BaseUrl = baseUrlP.String()

	reqBytes, err := proto.Marshal(&req)
	if err != nil {
		return nil, nil, err
	}

	var outSize uintptr

	reqSize := uintptr(len(reqBytes))

	coutSize := unsafe.Pointer(&outSize)
	creqSize := unsafe.Pointer(&reqSize)
	reqPtr := unsafe.Pointer(&reqBytes[0])

	cout := C.filter((*C.AdblockEngine)(engine.engine), (*C.uchar)(reqPtr), (*C.size_t)(creqSize), (*C.size_t)(coutSize))

	if outSize > 0 {
		defer C.drop_bytesarray(cout)
	}

	var resp FilterResponse

	protoBuf := unsafe.Slice((*byte)(cout), outSize)

	err = proto.Unmarshal(protoBuf, &resp)
	if err != nil {
		return nil, nil, err
	}

	return resp.Matches, resp.BlockedDomains, nil
}

func (engine *Engine) IdentifyBadElements(z *html.Node, baseUrl string) ([]string, []string, int, bool, error) {

	ids := []string{}
	classes := []string{}
	urls := []string{}

	containsGoogleTagManager := false

	walkHtmlNodesBadClasses(z, &ids, &classes, &urls, &containsGoogleTagManager)

	badIdsAndClasses, badUrls, err := engine.adblock.Filter(ids, classes, urls, baseUrl)
	if err != nil {
		return nil, nil, 0, false, err
	}

	badElementCount := len(badUrls)

	for _, ic := range badIdsAndClasses {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return nil, nil, 0, false, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			badElementCount++
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "block"})
		}
	}

	for _, ic := range badAreas {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return nil, nil, 0, false, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "skip"})
		}
	}

	for _, ic := range badClassesAndIds {
		sel, err := cascadia.Parse("#" + ic)
		if err != nil {
			return nil, nil, 0, false, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "skip"})
		}
	}

	for _, ic := range badClassesAndIds {
		sel, err := cascadia.Parse("." + ic)
		if err != nil {
			return nil, nil, 0, false, err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "skip"})
		}
	}

	return badIdsAndClasses, badUrls, badElementCount, containsGoogleTagManager, nil
}

func (engine *Engine) IdentifyGoodElements(z *html.Node, baseUrl string) error {

	for _, ic := range textTags {
		sel, err := cascadia.Parse(ic)
		if err != nil {
			return err
		}
		for _, a := range cascadia.QueryAll(z, sel) {
			a.Attr = append(a.Attr, html.Attribute{Key: "data-action", Val: "include"})
		}
	}

	return nil
}

func walkHtmlNodesBadClasses(n *html.Node, ids *[]string, classes *[]string, urls *[]string, containsGoogleTagManager *bool) {

	if n.Type == html.ElementNode {

		for _, attr := range n.Attr {
			if attr.Key == "class" {
				nc := strings.Split(attr.Val, " ")
				*classes = append(*classes, nc...)
			}
			if attr.Key == "id" {
				*ids = append(*classes, attr.Val)
			}
			if attr.Key == "href" {
				*urls = append(*urls, attr.Val)
			}
			if attr.Key == "src" {
				*urls = append(*urls, attr.Val)

				if strings.Contains(attr.Val, "googletagmanager") {
					*containsGoogleTagManager = true
				}
			}

		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkHtmlNodesBadClasses(c, ids, classes, urls, containsGoogleTagManager)
	}
}
