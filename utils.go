package main

import (
	"bytes"
	"net/url"
	"strings"
)

// regRemoveHtm := regexp.MustCompile(`(?m)(\s+)`)

func URLToPath(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)

	if err != nil {
		return "", err
	}

	var b bytes.Buffer

	b.WriteString(u.Host)

	path := strings.TrimRight(u.Path, "/")

	path = strings.ReplaceAll(path, ".htm", "")

	b.WriteString(path)

	if !strings.HasSuffix(path, ".html") {
		b.WriteString(".html")
	}

	return b.String(), nil
}

func inSliceInt(n int, h []int) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}

func chunkSlice(slice []int, chunkSize int) [][]int {
	var chunks [][]int
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func chunkHNUrlToCrawl(slice []*HNUrlToCrawl, chunkSize int) [][]*HNUrlToCrawl {
	var chunks [][]*HNUrlToCrawl
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func chunkFeedsToRefresh(slice []*Feed, chunkSize int) [][]*Feed {
	var chunks [][]*Feed
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func chunkArticlesToRefresh(slice []*Article, chunkSize int) [][]*Article {
	var chunks [][]*Article
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
