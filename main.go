package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/phea/greatfleas/article"
	"github.com/phea/greatfleas/core"
	"github.com/phea/greatfleas/router"

	"github.com/zenazn/goji"
	gojiweb "github.com/zenazn/goji/web"
)

func main() {
	log.Println("Starting blog application...")

	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	blog, err := core.NewBlogWithConfig(*configFile)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}
	blog.ParseTemplates()

	// Load articles
	if err := article.Parse(blog.ArticlePath); err != nil {
		log.Fatal(err)
		return
	}

	// Build feeds
	if err := core.BuildFeeds(article.ListAll()); err != nil {
		log.Fatal(err)
		return
	}

	// Setup static assets
	staticMux := gojiweb.New()
	staticMux.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(blog.PublicPath))))
	http.Handle("/assets/", staticMux)

	// Middleware
	goji.Use(blog.ApplyTemplates)

	// Routes
	goji.Get("/", blog.Route(router.Index))
	goji.Get("/articles", blog.Route(router.Articles))
	goji.Get("/articles/:num", blog.Route(router.ArticlesByPage))
	goji.Get("/article/:slug", blog.Route(router.Article))
	goji.Get("/archives", blog.Route(router.Archives))
	goji.Get("/about", blog.Route(router.About))
	goji.Get("/feed", blog.Route(router.Feed))

	flag.Set("bind", blog.Port)
	goji.Serve()
}
