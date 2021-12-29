package indecks_test

import (
	"context"
	"indecks"
	"os"
	"path"
	"testing"
)

func TestIndexer_NoProperty_Search(t *testing.T) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ix := &indecks.Indexer{
		Recursive:       true,
		Verbose:         true,
		ExtensionFilter: nil,
		Properties:      nil,
	}
	ch := ix.Search(ctx, path.Dir(os.Getenv("PWD")))
	var ies []*indecks.IndexEntry
	for ie := range ch {
		ies = append(ies, ie)
	}
	if len(ies) == 0 {
		t.Fatalf("no entries found")
	}
	root := "/Users/robgilham/go/src/indecks/"
	fn := path.Join(root, "indexer_test.go")
	if keyIndex(fn, ies) < 0 {
		t.Fatalf("expected to find %s, not found", fn)
	}
	fn = path.Join(root, "testdata/test.yaml")
	if keyIndex(fn, ies) < 0 {
		t.Fatalf("expected to find %s, not found", fn)
	}
}

func TestIndexer_YAMLSingleProperty_Search(t *testing.T) {
	propertyName := "testobject.one"
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ix := &indecks.Indexer{
		Recursive:       true,
		Verbose:         false,
		ExtensionFilter: map[string]bool{"yaml": true},
		Properties:      []indecks.PropertyParser{&indecks.YamlProperty{Names: []string{propertyName}}},
	}
	ch := ix.Search(ctx, os.Getenv("PWD"))
	var ies []*indecks.IndexEntry
	for ie := range ch {
		ies = append(ies, ie)
	}
	if len(ies) == 0 {
		t.Fatalf("no entries found")
	}
	root := "/Users/robgilham/go/src/indecks/"
	fn := path.Join(root, "testdata/test.yaml")
	ti := keyIndex(fn, ies)
	if ti < 0 {
		t.Fatalf("expected to find %s, not found", fn)
	}
	ie := ies[ti]
	if len(ie.Properties) != 1 {
		t.Fatalf("Incorrect property count. Expected %d, found %d", 1, len(ie.Properties))
	}
	p, ok := ie.Properties[propertyName]
	if !ok {
		t.Fatalf("failed to find test property %s", propertyName)
	}
	if p.String() != "1" {
		t.Fatalf("unexpected test property value for %s. Expected %s  Found %s", propertyName, "1", p.String())
	}
}

func TestIndexer_YAMLAllProperty_Search(t *testing.T) {
	propertyName := "testobject.*"
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ix := &indecks.Indexer{
		Recursive:       true,
		Verbose:         false,
		ExtensionFilter: map[string]bool{"yaml": true},
		Properties:      []indecks.PropertyParser{&indecks.YamlProperty{Names: []string{propertyName}}},
	}
	ch := ix.Search(ctx, os.Getenv("PWD"))
	var ies []*indecks.IndexEntry
	for ie := range ch {
		ies = append(ies, ie)
	}
	if len(ies) == 0 {
		t.Fatalf("no entries found")
	}
	root := "/Users/robgilham/go/src/indecks/"
	fn := path.Join(root, "testdata/test.yaml")
	ti := keyIndex(fn, ies)
	if ti < 0 {
		t.Fatalf("expected to find %s, not found", fn)
	}

	ie := ies[ti]
	if len(ie.Properties) != 3 {
		t.Fatalf("Incorrect property count. Expected %d, found %d", 3, len(ie.Properties))
	}

	p, ok := ie.Properties["testobject.one"]
	if !ok {
		t.Fatalf("failed to find test property %s", propertyName+".one")
	}
	if p.String() != "1" {
		t.Fatalf("unexpected test property value for %s. Expected %s  Found %s", propertyName, "1", p.String())
	}
	p, ok = ie.Properties["testobject.two"]
	if !ok {
		t.Fatalf("failed to find test property %s", propertyName+".two")
	}
	if p.String() != "2" {
		t.Fatalf("unexpected test property value for %s. Expected %s  Found %s", propertyName+".two", "2", p.String())
	}
	p, ok = ie.Properties["testobject.three"]
	if !ok {
		t.Fatalf("failed to find test property %s", propertyName+".three")
	}
	if p.String() != "3" {
		t.Fatalf("unexpected test property value for %s. Expected %s  Found %s", propertyName+".three", "3", p.String())
	}
}

func keyIndex(key string, ies []*indecks.IndexEntry) int {
	for i, ie := range ies {
		if key == ie.Key {
			return i
		}
	}
	return -1
}
