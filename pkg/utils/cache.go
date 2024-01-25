package utils

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/peterbourgon/diskv/v3"
)

type Cache struct {
	d *diskv.Diskv
}

func NewDiskCache(cachePath string) *Cache {
	d := diskv.New(diskv.Options{
		BasePath: cachePath,
		// CacheSizeMax:      1024 * 1024,
		CacheSizeMax:      10_737_418_240, // 1 GB
		AdvancedTransform: AdvancedTransformExample,
		InverseTransform:  InverseTransformExample,
	})

	return &Cache{d: d}
}

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

func (c *Cache) getHTMLKey(in string) (string, error) {
	u, err := url.Parse(in)
	if err != nil {
		return "", err
	}

	keyHash := fnv.New64()

	keyHash.Write([]byte(in))
	return fmt.Sprintf("%s/%v", u.Hostname(), keyHash.Sum64()), nil
}

func (c *Cache) ReadStream(in string) (io.ReadCloser, error) {
	k, err := c.getHTMLKey(in)
	if err != nil {
		return nil, err
	}

	htmlStream, err := c.d.ReadStream(k, false)
	if err != nil {
		return nil, err
	}
	return htmlStream, nil
}

func (c *Cache) WriteStream(in string, s io.Reader) error {
	k, err := c.getHTMLKey(in)
	if err != nil {
		return err
	}

	err = c.d.WriteStream(k, s, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(url string, client *http.Client) (io.Reader, error) {
	// check disk first
	disk, _ := c.ReadStream(url)
	if disk != nil {
		return disk, nil
	}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var buf bytes.Buffer
	tee := io.TeeReader(res.Body, &buf)
	err = c.WriteStream(url, tee)
	if err != nil {
		return nil, err
	}

	return tee, nil

}
