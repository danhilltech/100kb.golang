package http

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"

	"github.com/peterbourgon/diskv/v3"
)

func AdvancedTransformExample(key string) *diskv.PathKey {
	path := strings.Split(key, "/")
	last := len(path) - 1

	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last] + ".txt",
	}
}

// If you provide an AdvancedTransform, you must also provide its
// inverse:

func InverseTransformExample(pathKey *diskv.PathKey) (key string) {
	txt := pathKey.FileName[len(pathKey.FileName)-4:]
	if txt != ".txt" {
		panic("Invalid file found in storage folder!")
	}
	return strings.Join(pathKey.Path, "/") + pathKey.FileName[:len(pathKey.FileName)-4]
}

func getHTMLKey(r *http.Request) (string, error) {

	keyHash := fnv.New64()

	u := r.URL

	keyHash.Write([]byte(u.String()))
	return fmt.Sprintf("%s/%s-%v", u.Hostname(), r.Method, keyHash.Sum64()), nil
}

// func (c *Cache) ReadStream(in string) (io.ReadCloser, error) {
// 	k, err := c.getHTMLKey(in)
// 	if err != nil {
// 		return nil, err
// 	}

// 	htmlStream, err := c.d.ReadStream(k, false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return htmlStream, nil
// }

// func (c *Cache) WriteStream(in string, s io.Reader) error {
// 	k, err := c.getHTMLKey(in)
// 	if err != nil {
// 		return err
// 	}

// 	err = c.d.WriteStream(k, s, true)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Cache) Get(url string, client *http.Client) ([]byte, error) {
// 	// check disk first
// 	disk, _ := c.ReadStream(url)
// 	if disk != nil {
// 		byt, err := io.ReadAll(disk)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return byt, nil
// 	}

// 	res, err := client.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	byt, err := io.ReadAll(io.LimitReader(res.Body, 1_000_000))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// check for large files
// 	if len(byt) >= 1_000_000 {
// 		fmt.Printf("%s: %0.2f", url, float64(len(byt))/(1<<10))
// 	}

// 	io.Copy(io.Discard, res.Body)

// 	diskW := bytes.NewReader(byt)

// 	err = c.WriteStream(url, diskW)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return byt, nil

// }
