package tmpl

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

var defaultTemplateStore = instance()

var PRODUCTION = true

func Render(w http.ResponseWriter, name string, data Model) {
	defaultTemplateStore.render(w, name, data)
}

func ContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
}

func Reload() {
	defaultTemplateStore.load()
}

type Model map[string]interface{}

type templateStore struct {
	templates map[string]*template.Template
	bufpool   *bufferPool
	//funcs       template.FuncMap
	sync.RWMutex
}

func instance() *templateStore {
	t := &templateStore{
		templates: make(map[string]*template.Template),
		bufpool: &bufferPool{
			ch: make(chan *bytes.Buffer, 64),
		},
		//funcs:       defaultFuncMap,
	}
	t.load()
	return t
}

func (ts *templateStore) load() {
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

func (ts *templateStore) render(w http.ResponseWriter, name string, data Model) {
	var tmpl *template.Template
	if PRODUCTION {
		var ok bool
		ts.RLock()
		tmpl, ok = ts.templates[name]
		ts.RUnlock()
		if !ok {
			http.Error(w, "404. Page not found", 404)
			return
		}
	} else {
		includes, err := filepath.Glob("templates/includes/*.tmpl")
		if err != nil {
			log.Fatal(err)
		}
		files := append(includes, "templates/layouts/"+name)
		tmpl = template.Must(template.ParseFiles(files...))
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
