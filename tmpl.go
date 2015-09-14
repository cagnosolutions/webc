package web

import (
	"html/template"
	"net/http"
	"strings"
	"sync"
)

type Model map[string]interface{}

type TemplateStore struct {
	dir   string
	base  *template.Template
	cache map[string]*template.Template
	funcs template.FuncMap
	sync.RWMutex
}

func TemplateStoreInstance(dir string) *TemplateStore {
	ts := &TemplateStore{
		dir:   dir,
		base:  template.New("base.html"),
		cache: make(map[string]*template.Template),
		funcs: template.FuncMap{
			"title": strings.Title,
			"safe":  func(html string) template.HTML { return template.HTML(html) },
			"add":   func(a, b int) int { return a + b },
			"sub":   func(a, b int) int { return a - b },
			"decr":  func(a int) int { return a - 1 },
			"incr":  func(a int) int { return a + 1 },
			"split": strings.Split,
		},
	}
	ts.base.Funcs(ts.funcs)
	return ts
}

func (ts *TemplateStore) Cache(name ...string) {
	ts.Lock()
	for i := 0; i < len(name); i++ {
		ts.cache[name[i]] = template.Must(ts.base.ParseFiles(ts.dir+"/base.html", ts.dir+"/"+name[i]))
	}
	ts.Unlock()
}

func (ts *TemplateStore) Render(w http.ResponseWriter, name string, model interface{}) {
	var t *template.Template
	ts.RLock()
	t, ok := ts.cache[name]
	ts.RUnlock()
	if !ok {
		t = template.Must(ts.base.ParseFiles(ts.dir+"/base.html", ts.dir+"/"+name))
	}
	t.Execute(w, model)
}

func ContentType(w http.ResponseWriter, typ string) {
	w.Header().Set("Content Type", typ)
}

/*
func (ts *TemplateEngine) ToString(name string, model interface{}) string {
	var buf bytes.Buffer
	ts.cache[name].Execute(&buf, model)
	return buf.String()
}

func (ts *TemplateEngine) ToBytes(name string, model interface{}) []byte {
	var buf bytes.Buffer
	ts.cache[name].Execute(&buf, model)
	return buf.Bytes()
}
*/
