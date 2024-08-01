package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

func (engine *Engine) getTopStories() ([]int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/topstories.json", HN_BASE))
	// handle the error if there is one
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var out []int

	err = json.Unmarshal(body, &out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (engine *Engine) getHNItem(id int) (*ToCrawl, error) {
	tmpItem := &ToCrawl{HNID: id}

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

	var item ToCrawl

	err = json.Unmarshal(body, &item)
	if err != nil {
		return tmpItem, err
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

func (engine *Engine) getHNWorker(jobs <-chan int, results chan<- *ToCrawl) {
	for id := range jobs {
		item, err := engine.getHNItem(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- item
	}
}

// Gets the latest content from Hacker news
func (engine *Engine) RunHNRefresh(ctx context.Context, chunkSize int, totalFetch int, workers int) error {

	max, err := engine.getMaxId()
	if err != nil {
		return err
	}

	min := max - totalFetch

	var ids []int

	existingIds, err := engine.getExistingIDs()
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

	topIds, err := engine.getTopStories()
	if err != nil {
		return err
	}
	ids = append(ids, topIds...)

	fmt.Printf("Getting %d HN items\n", len(ids))

	jobs := make(chan int, len(ids))
	results := make(chan *ToCrawl, len(ids))

	for w := 1; w <= workers; w++ {
		go engine.getHNWorker(jobs, results)
	}

	for j := 1; j <= len(ids); j++ {
		jobs <- ids[j-1]
	}
	close(jobs)

	t := time.Now().UnixMilli()

	txn, _ := engine.db.Begin()
	defer txn.Rollback()

	for a := 1; a <= len(ids); a++ {
		select {
		case <-ctx.Done():
			txn.Commit()
			return ctx.Err()
		case item := <-results:
			err = engine.InsertToCrawl(txn, item)
			if err != nil {
				return err
			}

			if a > 0 && a%chunkSize == 0 {
				err := txn.Commit()
				if err != nil {
					return err
				}
				txn, _ = engine.db.Begin()

				diff := time.Now().UnixMilli() - t
				qps := (float64(chunkSize) / float64(diff)) * 1000
				t = time.Now().UnixMilli()
				fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(ids), qps)

			}

		}
	}
	txn.Commit()
	fmt.Printf("\tdone %d\n", len(ids))

	return nil
}
