package handler

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

type Renderer interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type TemplateMap struct {
	templates map[string]*template.Template
}

func NewTemplateMap(dir string, funcMap template.FuncMap) (*TemplateMap, error) {
	base := filepath.Join(dir, "base.html")
	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}

	tm := &TemplateMap{templates: make(map[string]*template.Template)}
	for _, page := range pages {
		name := filepath.Base(page)
		if name == "base.html" {
			continue
		}
		t, err := template.New(name).Funcs(funcMap).ParseFiles(base, page)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
		tm.templates[name] = t
	}
	return tm, nil
}

func (tm *TemplateMap) ExecuteTemplate(wr io.Writer, name string, data any) error {
	t, ok := tm.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}
	return t.ExecuteTemplate(wr, "base", data)
}
