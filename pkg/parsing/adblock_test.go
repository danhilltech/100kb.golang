package parsing

import (
	"bufio"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/html"
)

func TestAdBlock(t *testing.T) {

	adblock, err := NewAdblockEngine()
	if err != nil {
		t.Fatal(err)
	}

	files := []string{}

	root := "/workspaces/100kb.golang/.cache"
	fileSystem := os.DirFS(root)

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return nil
		}

		if strings.Contains(path, "GET-") && len(files) < 1000 {
			files = append(files, path)
		}

		return nil
	})

	// files, err := filepath.Glob("/workspaces/100kb.golang/**/GET-*.txt")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	for i := range files {
		j := rand.Intn(i + 1)
		files[i], files[j] = files[j], files[i]
	}

	fmt.Println(files[:100])

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		// time.Sleep(5 * time.Millisecond)
		go func(f string) {
			defer wg.Done()
			err := fakeParse(adblock, f)
			if err != nil {
				fmt.Println(err)
			}
		}(files[i])
	}
	wg.Wait()

	defer adblock.Close()

}

func fakeParse(adblock *AdblockEngine, p string) error {
	parseAnalysis := ParseAnalysis{
		Ids:           make([]string, 0),
		Classes:       make([]string, 0),
		Urls:          make([]string, 0),
		Links:         make([]string, 0),
		BadUrls:       make([]string, 0),
		BadElements:   make([]string, 0),
		BadLinkTitles: make([]string, 0),
	}

	f, err := os.Open(fmt.Sprintf("/workspaces/100kb.golang/.cache/%s", p))
	if err != nil {
		return err
	}

	buf := bufio.NewReader(f)

	cached, err := http.ReadResponse(buf, nil)
	if err != nil {
		return err
	}
	defer cached.Body.Close()

	n, err := html.Parse(cached.Body)
	if err != nil {
		return err
	}

	walkHtmlNodesAndIdentify(n, &parseAnalysis)

	// parseAnalysis.Ids = append(parseAnalysis.Ids, "google_ads")
	// parseAnalysis.Classes = append(parseAnalysis.Classes, "google_ads")
	// parseAnalysis.Urls = append(parseAnalysis.Urls, "https://www.googletagmanager.com")
	_, _, err = adblock.Filter(parseAnalysis.Ids, parseAnalysis.Classes, parseAnalysis.Urls, "https://www.danhill.is")

	// fmt.Println(ids, classes)
	return err
}
