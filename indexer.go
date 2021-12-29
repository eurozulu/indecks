package indecks

import (
	"context"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Indexer struct {
	Recursive       bool
	Verbose         bool
	ExtensionFilter map[string]bool
	Properties      []PropertyParser
}

func (i Indexer) Search(ctx context.Context, p ...string) <-chan *IndexEntry {
	ch := make(chan *IndexEntry)
	go func(ch chan<- *IndexEntry) {
		defer close(ch)
		var wg sync.WaitGroup
		for _, fp := range p {
			wg.Add(1)
			go i.search(ctx, fp, ch, &wg)
		}
		wg.Wait()
	}(ch)
	return ch
}

func (i Indexer) search(ctx context.Context, root string, out chan<- *IndexEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	var wwg sync.WaitGroup
	// Kick off a routine to load each eligible file into a indexEntry
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if !i.handleError(err) {
			return nil
		}
		if strings.HasPrefix(path.Base(p), ".") && p != root {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if i.ExtensionFilter != nil {
			if !i.ExtensionFilter[strings.TrimLeft(path.Ext(p), ".")] {
				return nil
			}
		}
		if d.IsDir() {
			if !i.Recursive && p != root {
				return filepath.SkipDir
			}
			return nil
		}
		// Read and parse file independantly
		wwg.Add(1)
		go func(p string, t time.Time, wg *sync.WaitGroup) {
			defer wg.Done()
			by, err := ioutil.ReadFile(p)
			if !i.handleError(err) {
				return
			}
			props := i.parseProperties(p, by)
			select {
			case <-ctx.Done():
				return
			case out <- &IndexEntry{
				Key:        p,
				Properties: props,
			}:
			}
		}(p, time.Now(), &wwg)
		return nil
	})
	wwg.Wait()
	if !i.handleError(err) {
		return
	}
}

func (i Indexer) handleError(err error) bool {
	if err == nil {
		return true
	}
	if i.Verbose {
		log.Println(err)
	}
	return false
}

func (i Indexer) parseProperties(p string, by []byte) map[string]Property {
	var props map[string]Property
	for _, pp := range i.Properties {
		pm, err := pp.ParseProperties(p, by)
		if !i.handleError(err) {
			continue
		}
		props = appendMap(props, pm)
	}
	return props
}

func appendMap(m map[string]Property, ms map[string]Property) map[string]Property {
	if len(ms) == 0 {
		return m
	}
	if m == nil {
		m = map[string]Property{}
	}
	for k, v := range ms {
		m[k] = v
	}
	return m
}
