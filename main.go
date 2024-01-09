package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting")

	err := initDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stopDB()

	// err = RunHackerNewsRefresh()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = RunNewFeedSearch()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = RunFeedRefresh()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = RunArticleRefresh()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	ai, err := loadAi("models/bert3.bin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ai.Close()

	fmt.Println(ai)

	vec, err := ai.Embeddings([]string{"My name is dan"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(vec)

	// err = RunArticleMeta()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	dbTidy()

}
