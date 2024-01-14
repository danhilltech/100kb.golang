package crawler

func Chunk(slice []*UrlToCrawl, chunkSize int) [][]*UrlToCrawl {
	var chunks [][]*UrlToCrawl
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
