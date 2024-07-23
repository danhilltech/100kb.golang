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
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

type AdblockEngine struct {
	engine *C.AdblockEngine
	mutex  sync.Mutex
	pinner runtime.Pinner
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

	engine := AdblockEngine{}
	engine.pinner = runtime.Pinner{}

	engine.pinner.Pin(reqBytes)

	reqSize := uintptr(len(reqBytes))

	creqSize := unsafe.Pointer(&reqSize)
	reqPtr := unsafe.Pointer(&reqBytes[0])

	a := C.new_adblock((*C.uchar)(reqPtr), (*C.size_t)(creqSize))
	engine.pinner.Pin(a)
	engine.engine = a

	return &engine, nil

}

func (engine *AdblockEngine) Close() {
	C.drop_adblock((*C.AdblockEngine)(engine.engine))
	engine.pinner.Unpin()
}

func (engine *AdblockEngine) Filter(ids []string, classes []string, urls []string, baseUrl string) ([]string, []string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
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
