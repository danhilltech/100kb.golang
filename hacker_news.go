package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const HN_BASE = "https://hacker-news.firebaseio.com/v0"

var hnClient *http.Client

type HNItemType string

type HNItem struct {
	ID    int
	URL   string
	By    string
	Type  HNItemType
	Time  int
	Score int
}

const (
	HNItemTypeStory = "story"
)

func init() {
	tr := &http.Transport{MaxIdleConnsPerHost: 1024, TLSHandshakeTimeout: 0 * time.Second}
	hnClient = &http.Client{Transport: tr}
}

func getMaxId() (int, error) {
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

func getItem(id int) (*HNItem, error) {
	resp, err := hnClient.Get(fmt.Sprintf("%s/item/%d.json", HN_BASE, id))

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

func getItemWorker(jobs <-chan int, results chan<- *HNItem) {
	for id := range jobs {
		item, err := getItem(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- item
	}
}

func getItems(ids []int) ([]*HNItem, error) {

	workers := 10
	jobs := make(chan int, len(ids))
	results := make(chan *HNItem, len(ids))

	for w := 1; w <= workers; w++ {
		go getItemWorker(jobs, results)
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
