package utils

import (
	"bytes"
	"database/sql"
	"net/url"
	"strings"
)

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

func InSliceInt(n int, h []int) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}

func ChunkSlice(slice []int, chunkSize int) [][]int {
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

func NullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NullInt64(s int64) sql.NullInt64 {
	if s == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: s,
		Valid: true,
	}
}

func NullFloat64(s float64) sql.NullFloat64 {
	if s == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: s,
		Valid:   true,
	}
}
