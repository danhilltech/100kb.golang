package parsing

import (
	"bufio"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func TestAdBlock(t *testing.T) {

	numToRun := 20000
	threads := runtime.NumCPU()

	adblock, err := NewAdblockEngine(nil)
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

		if strings.Contains(path, "GET-") && len(files) < numToRun {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// files, err := filepath.Glob("/workspaces/100kb.golang/**/GET-*.txt")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// for i := range files {
	// 	j := rand.Intn(i + 1)
	// 	files[i], files[j] = files[j], files[i]
	// }

	fmt.Println("starting...")

	var wg sync.WaitGroup

	size := 0

	ticker := time.NewTicker(150 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				runtime.GC()
			}
		}
	}()

	for i := 0; i < numToRun; i++ {
		wg.Add(1)

		// time.Sleep(5 * time.Millisecond)
		go func(f string) {
			defer wg.Done()
			err := fakeParse(adblock, f)
			if err != nil {
				fmt.Println(err)
			}
		}(files[i])

		if size%threads == 0 {
			wg.Wait()

		}
		if size%20 == 0 {
			fmt.Printf(".")
		}

		size++
	}

	wg.Wait()
	ticker.Stop()
	done <- true

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

	// rd := strings.NewReader(`garb`)

	// n, err := html.Parse(rd)
	// if err != nil {
	// 	return err
	// }

	// engine.log.Println(p)

	walkHtmlNodesAndIdentify(n, &parseAnalysis)

	// parseAnalysis.Ids = append(parseAnalysis.Ids, "google_ads")
	// parseAnalysis.Classes = append(parseAnalysis.Classes, "google_ads")
	// parseAnalysis.Urls = append(parseAnalysis.Urls, "https://www.googletagmanager.com")
	_, _, err = adblock.Filter(parseAnalysis.Ids, parseAnalysis.Classes, parseAnalysis.Urls, "https://www.danhill.is")

	// engine.log.Println(ids, classes)
	return err
}
