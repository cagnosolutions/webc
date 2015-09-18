package tmpl

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type Model map[string]interface{}

type TemplateStore struct {
	templates   map[string]*template.Template
	bufpool     *bufferPool
	Development bool
	//funcs       template.FuncMap
	sync.RWMutex
}

func NewTemplateStore() *TemplateStore {
	t := &TemplateStore{
		templates: make(map[string]*template.Template),
		bufpool: &bufferPool{
			ch: make(chan *bytes.Buffer, 64),
		},
		//funcs:       defaultFuncMap,
	}
	t.Load()
	return t
}

func (ts *TemplateStore) Load() {
	layouts, err := filepath.Glob("templates/layouts/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob("templates/includes/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		ts.Lock()
		ts.templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
		ts.Unlock()
	}

}

func (ts *TemplateStore) Render(w http.ResponseWriter, name string, data Model) {
	var tmpl *template.Template
	if ts.Development {
		includes, err := filepath.Glob("templates/includes/*.tmpl")
		if err != nil {
			log.Fatal(err)
		}
		files := append(includes, "templates/layouts/"+name)
		tmpl = template.Must(template.ParseFiles(files...))
	} else {
		var ok bool
		ts.RLock()
		tmpl, ok = ts.templates[name]
		ts.RUnlock()
		if !ok {
			http.Error(w, "404. Page not found", 404)
			return
		}
	}

	buf := ts.bufpool.get()
	defer ts.bufpool.reset(buf)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "500. Internal server error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return
}

func ContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
}
