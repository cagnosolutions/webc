package tmpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cagnosolutions/webc/util"
)

type Model map[string]interface{}

type TemplateStore struct {
	templates   map[string]*template.Template
	bufpool     *bufferPool
	Development bool
	TemplateDir string
	funcs       template.FuncMap
	sync.RWMutex
}

func NewTemplateStore(development bool) *TemplateStore {
	t := &TemplateStore{
		templates: make(map[string]*template.Template),
		bufpool: &bufferPool{
			ch: make(chan *bytes.Buffer, 64),
		},
		TemplateDir: "templates/",
		Development: development,
		funcs: template.FuncMap{
			"title": strings.Title,
			"safe":  func(html string) template.HTML { return template.HTML(html) },
			"add":   func(a, b int) int { return a + b },
			"sub":   func(a, b int) int { return a - b },
			"decr":  func(a int) int { return a - 1 },
			"incr":  func(a int) int { return a + 1 },
			"split": strings.Split,
			"map":   func(a map[string]string, b string) interface{} { return a[b] },
			"pretty": func(v interface{}) string {
				b, err := json.MarshalIndent(v, "", "\t")
				if err != nil {
					log.Println(err)
				}
				return string(b)
			},
			"date": func(d string) string {
				ds := util.SliceString(d, '-')
				return fmt.Sprintf("%s/%s/%s", ds[1], ds[2], ds[0])
			},
		},
	}
	t.Load()
	return t
}

func (ts *TemplateStore) Load() {
	layouts, err := filepath.Glob(ts.TemplateDir + "layouts/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(ts.TemplateDir + "includes/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	for _, layout := range layouts {
		files := append(includes, layout)
		ts.Lock()
		ts.templates[filepath.Base(layout)] = template.Must(template.New("func").Funcs(ts.funcs).ParseFiles(files...))
		ts.Unlock()
	}

}

func (ts *TemplateStore) Render(w http.ResponseWriter, name string, data Model) {
	var tmpl *template.Template
	if ts.Development {
		includes, err := filepath.Glob(ts.TemplateDir + "includes/*.tmpl")
		if err != nil {
			log.Fatal(err)
		}
		files := append(includes, ts.TemplateDir+"layouts/"+name)
		tmpl = template.Must(template.New("func").Funcs(ts.funcs).ParseFiles(files...))
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
