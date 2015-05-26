// Package article
package article

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

const (
	dateFormat     = "2006-01-02"
	resultsPerPage = 10
)

// db is the package wide object that holds the main Articles struct.
var db *Articles

// An Entry holds a single blog entry.
type Entry struct {
	Title   string
	Slug    string
	Content template.HTML
	Created time.Time
	Updated time.Time
}

type ArticleIndex []*Entry

// Articles is the object that holds all the Entries.
type Articles struct {
	Index ArticleIndex
	Table map[string]*Entry

	total int
}

// Len
func (idx ArticleIndex) Len() int {
	return len(idx)
}

// Less sorts so the newest articles are first.
func (idx ArticleIndex) Less(i, j int) bool {
	return idx[j].Created.Before(idx[i].Created)
}

// Swap
func (idx ArticleIndex) Swap(i, j int) {
	idx[i], idx[j] = idx[j], idx[i]
}

// Get
func Get(slug string) (*Entry, error) {
	e, ok := db.Table[slug]
	if !ok {
		return nil, errors.New("article not found.")
	}

	return e, nil
}

// List
func List(n int) ([]*Entry, error) {
	if db == nil {
		return nil, errors.New("articles not found.")
	}

	if db.total <= 5 {
		return db.Index, nil
	}

	low := (n - 1) * resultsPerPage
	high := n * resultsPerPage

	return db.Index[low : high-1], nil
}

// ListAll
func ListAll() []*Entry {
	return db.Index
}

// Recent
func Recent(n int) ([]*Entry, error) {
	if db == nil {
		return nil, errors.New("articles not found.")
	}

	if db.total <= n {
		return db.Index, nil
	}

	return db.Index[:n], nil
}

func Parse(path string) error {
	idx, err := parse(path)
	if err != nil {
		return err
	}

	articles := make(map[string]*Entry, len(idx))
	for _, i := range idx {
		articles[i.Slug] = i
	}

	sort.Sort(idx)
	db = &Articles{idx, articles, len(idx)}
	return nil
}

// parse walks through the provided directory and parses each .article file.
func parse(dir string) (ArticleIndex, error) {
	var list ArticleIndex
	// Walk directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".md" {
			//return filepath.SkipDir
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		e, err := parseFile(data)
		if err != nil {
			return err
		}

		list = append(list, e)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

func parseFile(data []byte) (*Entry, error) {
	e := &Entry{}
	var err error
	buf := bytes.NewBuffer(data)
	e.Title, _ = buf.ReadString('\n')
	e.Slug, _ = buf.ReadString('\n')
	created, _ := buf.ReadString('\n')
	updated, _ := buf.ReadString('\n')
	e.Content = template.HTML(blackfriday.MarkdownCommon(buf.Bytes()))

	e.Created, err = time.Parse(dateFormat, strings.TrimSpace(created))
	if err != nil {
		return nil, err
	}
	e.Updated, err = time.Parse(dateFormat, strings.TrimSpace(updated))
	if err != nil {
		return nil, err
	}

	e.Title = strings.TrimSpace(e.Title)
	e.Slug = strings.TrimSpace(e.Slug)
	return e, nil
}
