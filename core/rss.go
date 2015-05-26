package core

import (
	"io/ioutil"
	"time"

	"github.com/gorilla/feeds"
	"github.com/phea/greatfleas/article"
)

func BuildFeeds(list []*article.Entry) error {
	feed := newFeed()

	var items []*feeds.Item
	for _, a := range list {
		item := &feeds.Item{
			Title:       a.Title,
			Link:        &feeds.Link{Href: "http://www.greatfleas.com/" + a.Slug},
			Description: string(a.Content),
			Author:      &feeds.Author{"Phea Duch", "email@email.com"},
			Created:     a.Created,
		}

		items = append(items, item)
	}

	feed.Items = items
	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	// Write to file
	ioutil.WriteFile("public/atom.xml", []byte(atom), 0644)
	ioutil.WriteFile("public/rss.xml", []byte(rss), 0644)

	return nil
}

func newFeed() *feeds.Feed {
	return &feeds.Feed{
		Title:       "GreatFleas",
		Link:        &feeds.Link{Href: "http://www.greatfleas.com"},
		Description: "Thoughts from a founder.",
		Author:      &feeds.Author{"Phea Duch", "email@email.com"},
		Created:     time.Now(),
	}
}
