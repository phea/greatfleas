// Package router
package router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/phea/greatfleas/article"
	"github.com/phea/greatfleas/core"
	"github.com/zenazn/goji/web"
)

// Home
func Index(c web.C, r *http.Request) (string, int) {
	t := core.GetTemplate(c)
	list, err := article.List(0)
	if err != nil {
		return fmt.Sprintf("%s", err), http.StatusInternalServerError
	}
	recentList, _ := article.Recent(5)

	c.Env["Articles"] = list
	c.Env["RecentArticles"] = recentList
	articles := core.Render(t, "home", c.Env)

	c.Env["Content"] = template.HTML(articles)
	return core.Render(t, "main", c.Env), http.StatusOK
}

// Error (temp)
func Error(c web.C, r *http.Request) (string, int) {
	return "This is an error", http.StatusInternalServerError
}

// Articles handles /articles route.
func Articles(c web.C, r *http.Request) (string, int) {
	t := core.GetTemplate(c)
	list, err := article.List(0)
	if err != nil {
		return fmt.Sprintf("%s", err), http.StatusInternalServerError
	}

	recentList, _ := article.Recent(5)

	c.Env["Articles"] = list
	c.Env["RecentArticles"] = recentList
	articles := core.Render(t, "articles", c.Env)

	c.Env["Content"] = template.HTML(articles)
	return core.Render(t, "main", c.Env), http.StatusOK
}

// ArticlesByPage handles /articles/:num route.
func ArticlesByPage(c web.C, r *http.Request) (string, int) {
	t := core.GetTemplate(c)
	num := c.URLParams["num"]

	n, _ := strconv.Atoi(num)
	list, err := article.List(n)
	if err != nil {
		return fmt.Sprintf("%s", err), http.StatusInternalServerError
	}
	recentList, _ := article.Recent(5)

	c.Env["Articles"] = list
	c.Env["RecentArticles"] = recentList
	articles := core.Render(t, "articles", c.Env)

	c.Env["Content"] = template.HTML(articles)
	return core.Render(t, "main", c.Env), http.StatusOK
}

// Article handles /article/:slug route.
func Article(c web.C, r *http.Request) (string, int) {
	t := core.GetTemplate(c)
	slug := c.URLParams["slug"]
	a, err := article.Get(slug)
	if err != nil {
		return "Article not found.", http.StatusNotFound
	}
	recentList, _ := article.Recent(5)

	c.Env["Article"] = a
	c.Env["RecentArticles"] = recentList
	article := core.Render(t, "article", c.Env)

	c.Env["Content"] = template.HTML(article)
	return core.Render(t, "main", c.Env), http.StatusOK
}

// Archives handles /archives route.
func Archives(c web.C, r *http.Request) (string, int) {
	t := core.GetTemplate(c)
	list, err := article.List(0)
	if err != nil {
		return fmt.Sprintf("%s", err), http.StatusInternalServerError
	}
	recentList, _ := article.Recent(5)

	c.Env["Articles"] = list
	c.Env["RecentArticles"] = recentList
	archives := core.Render(t, "archives", c.Env)

	c.Env["Content"] = template.HTML(archives)
	return core.Render(t, "main", c.Env), http.StatusOK
}

// About handles /about route.
func About(c web.C, r *http.Request) (string, int) {
	t := core.GetTemplate(c)
	recentList, _ := article.Recent(5)

	c.Env["RecentArticles"] = recentList
	about := core.Render(t, "about", c.Env)

	c.Env["Content"] = template.HTML(about)
	return core.Render(t, "main", c.Env), http.StatusOK
}

// Feed handles /feed route, it returns an atom.xml file.
func Feed(c web.C, r *http.Request) (string, int) {
	xml, err := ioutil.ReadFile("public/atom.xml")
	if err != nil {
		return "XML Feed Unavailable", http.StatusOK
	}

	return string(xml), http.StatusOK
}
