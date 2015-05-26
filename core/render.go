package core

import (
	"bytes"
	"html/template"
)

// Render
func Render(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}
