package collection

import (
	"fmt"
	"io/fs"
	"maps"
	"sync"

	"github.com/Panterrich/PhotoStudio/pkg/image"
	"github.com/charlievieth/fastwalk"
)

type Image struct {
	Jpeg string // path to jpeg
	Raw  string // path to raw
}

type Images struct {
	Jpegs map[string]struct{}
	Raws  map[string]struct{}
}

type Collection struct {
	mutex  sync.RWMutex
	images Images
}

func NewCollection(srcDir string, nWorkers int) (*Collection, error) {
	c := &Collection{
		mutex: sync.RWMutex{},
		images: Images{
			Jpegs: make(map[string]struct{}),
			Raws:  make(map[string]struct{}),
		},
	}

	walkFn := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}

		if !info.IsDir() {
			c.AddImage(path)
		}

		return nil
	}

	cfg := &fastwalk.Config{
		Follow:     false,
		ToSlash:    false,
		Sort:       fastwalk.SortDirsFirst,
		NumWorkers: nWorkers,
	}

	if err := fastwalk.Walk(cfg, srcDir, walkFn); err != nil {
		return nil, fmt.Errorf("create new image collection error: %w", err)
	}

	return c, nil
}

func (c *Collection) AddImage(path string) {
	if image.IsJpeg(path) {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		c.images.Jpegs[path] = struct{}{}
	} else if image.IsRaw(path) {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		c.images.Raws[path] = struct{}{}
	}
}

func (c *Collection) Size() (int, int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.images.Jpegs), len(c.images.Raws)
}

func (c *Collection) Images() Images {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return Images{
		Jpegs: maps.Clone(c.images.Jpegs),
		Raws:  maps.Clone(c.images.Raws),
	}
}
