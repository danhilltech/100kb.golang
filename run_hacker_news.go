package main

import "fmt"

func RunHackerNewsRefresh() error {
	max, err := getMaxId()
	if err != nil {
		return err
	}

	chunkSize := 200
	totalFetch := 20_000

	min := max - totalFetch

	var ids []int

	existingIds, err := getExistingHNIDs()
	if err != nil {
		return err
	}

	for i := min; i < max; i++ {
		if inSliceInt(i, existingIds) {
			continue
		}
		ids = append(ids, i)
	}

	chunkIds := chunkSlice(ids, chunkSize)

	for _, chunk := range chunkIds {
		err = doHackerNewsChunk(chunk)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func doHackerNewsChunk(chunk []int) error {
	fmt.Println("Starting chunk")
	items, err := getItems(chunk)
	if err != nil {
		return err
	}

	insertTxn, err := db.Begin()
	defer insertTxn.Rollback()

	if err != nil {
		return err
	}
	for _, item := range items {
		err = saveHNItem(item, insertTxn)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Saving to db")
	err = insertTxn.Commit()
	if err != nil {
		return err
	}
	return nil
}
