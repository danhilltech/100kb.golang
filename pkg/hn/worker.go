package hn

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"mvdan.cc/xurls/v2"
)

// 39,190,942

func (engine *Engine) getMaxId() (int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/maxitem.json", HN_BASE))
	// handle the error if there is one
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	var max int

	err = json.Unmarshal(body, &max)

	if err != nil {
		return 0, err
	}

	return max, nil
}

func (engine *Engine) getItem(id int) (*HNItem, error) {
	tmpItem := &HNItem{ID: id}

	resp, err := engine.client.Get(fmt.Sprintf("%s/item/%d.json", HN_BASE, id))

	// handle the error if there is one
	if err != nil {
		return tmpItem, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return tmpItem, err
	}

	var item HNItem

	err = json.Unmarshal(body, &item)
	if err != nil {
		return tmpItem, err
	}

	if item.Type == "comment" {
		rxStrict := xurls.Strict()

		txt := html.UnescapeString(item.Text)

		urls := rxStrict.FindAllString(txt, 1)
		if len(urls) >= 1 && strings.HasPrefix(urls[0], "https://") {
			item.URL = urls[0]
		}
	}

	if item.URL != "" {
		u, err := url.Parse(item.URL)
		if err != nil {
			item.URL = ""
		} else {
			item.Domain = u.Hostname()
		}
	}

	return &item, nil
}

func (engine *Engine) getItemWorker(jobs <-chan int, results chan<- *HNItem) {
	for id := range jobs {
		item, err := engine.getItem(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- item
	}
}

// Gets the latest content from Hacker news
func (engine *Engine) RunRefresh(chunkSize int, totalFetch int, workers int) error {

	max, err := engine.getMaxId()
	if err != nil {
		return err
	}

	min := max - totalFetch

	var ids []int

	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	existingIds, err := engine.getExistingIDs(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	wanted := map[int]bool{}

	for i := min; i < max; i++ {
		wanted[i] = true
	}

	for _, e := range existingIds {
		wanted[e] = false
	}

	for i := min; i < max; i++ {
		if wanted[i] {
			ids = append(ids, i)
		}
	}

	fmt.Printf("Getting %d HN items\n", len(ids))

	jobs := make(chan int, len(ids))
	results := make(chan *HNItem, len(ids))

	for w := 1; w <= workers; w++ {
		go engine.getItemWorker(jobs, results)
	}

	for j := 1; j <= len(ids); j++ {
		jobs <- ids[j-1]
	}
	close(jobs)

	insertTxn, err := engine.db.Begin()
	defer insertTxn.Rollback()
	if err != nil {
		return err
	}

	t := time.Now().UnixMilli()

	for a := 1; a <= len(ids); a++ {
		item := <-results

		err = engine.save(item, insertTxn)
		if err != nil {
			return err
		}

		if a > 0 && a%chunkSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(chunkSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(ids), qps)
			err = insertTxn.Commit()
			if err != nil {
				return err
			}
			insertTxn, err = engine.db.Begin()
			if err != nil {
				return err
			}
		}

	}
	fmt.Printf("\tdone %d\n", len(ids))

	err = insertTxn.Commit()
	if err != nil {
		return err
	}

	return nil
}
