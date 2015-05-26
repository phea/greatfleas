// Package core
package core

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zenazn/goji/web"
)

// Blog is the main application structure.
type Blog struct {
	*Config
	Template *template.Template
}

func NewBlogWithConfig(file string) (*Blog, error) {
	conf, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	return &Blog{conf, nil}, nil
}

var templateFuncs = template.FuncMap{
	"ts": func(ts time.Time) string {
		tString := ts.String()
		return tString[:10]
	},
}

func (b *Blog) ParseTemplates() error {
	var templates []string

	fn := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			templates = append(templates, path)
		}
		return nil
	}

	if err := filepath.Walk(b.TemplatePath, fn); err != nil {
		return err
	}

	tmpl := template.New("template")
	b.Template = template.Must(tmpl.Funcs(templateFuncs).ParseFiles(templates...))
	return nil
}

type Controller func(web.C, *http.Request) (string, int)

// Route is a wrapper for http requests.
func (b *Blog) Route(controller Controller) web.HandlerFunc {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		body, code := controller(c, r)

		switch code {
		case http.StatusOK:
			io.WriteString(w, body)
		case http.StatusInternalServerError:
			log.Println("Error: ", body)
			log.Printf("%+v\n", r)
			io.WriteString(w, "Error")
		}
	}
	return fn
}

// Helper section
// GetTemplate
func GetTemplate(c web.C) *template.Template {
	return c.Env["Template"].(*template.Template)
}
