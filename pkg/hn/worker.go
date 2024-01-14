package hn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

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
	resp, err := engine.client.Get(fmt.Sprintf("%s/item/%d.json", HN_BASE, id))

	// handle the error if there is one
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var item HNItem

	err = json.Unmarshal(body, &item)

	if err != nil {
		return nil, err
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

func (engine *Engine) getItems(ids []int, workers int) ([]*HNItem, error) {
	jobs := make(chan int, len(ids))
	results := make(chan *HNItem, len(ids))

	for w := 1; w <= workers; w++ {
		go engine.getItemWorker(jobs, results)
	}

	for j := 1; j <= len(ids); j++ {
		jobs <- ids[j-1]
	}
	close(jobs)

	items := make([]*HNItem, len(ids))

	for a := 1; a <= len(ids); a++ {
		b := <-results
		items[a-1] = b
	}

	return items, nil
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

	for i := min; i < max; i++ {
		if utils.InSliceInt(i, existingIds) {
			continue
		}
		ids = append(ids, i)
	}

	chunkIds := utils.ChunkSlice(ids, chunkSize)

	fmt.Printf("Getting %d HN items in %d chunks\n", len(ids), len(chunkIds))

	for _, chunk := range chunkIds {
		err = engine.doHackerNewsChunk(chunk, workers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (engine *Engine) doHackerNewsChunk(chunk []int, workers int) error {
	fmt.Printf("Chunk...\t\t")
	defer fmt.Printf("✅\n")
	items, err := engine.getItems(chunk, workers)
	if err != nil {
		return err
	}

	insertTxn, err := engine.db.Begin()
	defer insertTxn.Rollback()

	if err != nil {
		return err
	}
	for _, item := range items {
		err = engine.save(item, insertTxn)
		if err != nil {
			return err
		}
	}

	err = insertTxn.Commit()
	if err != nil {
		return err
	}
	return nil
}
