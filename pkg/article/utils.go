package article

import (
	"fmt"
	"hash/fnv"
	"net/url"
)

func Chunk(slice []*Article, chunkSize int) [][]*Article {
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

func (a *Article) getHTMLKey() (string, error) {
	u, err := url.Parse(a.Url)
	if err != nil {
		return "", err
	}

	keyHash := fnv.New64()

	k := u.EscapedPath() + u.Query().Encode()

	keyHash.Write([]byte(k))
	return fmt.Sprintf("%s/%v", u.Hostname(), keyHash.Sum64()), nil
}
