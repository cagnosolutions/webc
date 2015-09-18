package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type M map[string]interface{}

/*
type TemplateStore *template.Template

var T *template.Template

func RegisterTemplates(fileNames ...string) {
	T = template.Must(template.ParseFiles("templates/base.html"))
	T.ParseFiles(fileNames...)


	for _, t := range T.Templates() {
		fmt.Printf("%+v\n",t.Name())
	}
}
*/

type TemplateStore struct {
	dir   string
	base  *template.Template
	cache map[string]*template.Template
	funcs template.FuncMap
	sync.RWMutex
}

func TemplateStoreInstance(dir string, funcs map[string]func()) *TemplateStore {
	ts := &TemplateStore{
		dir: dir,
		//base:  template.New("base.html"),
		base:  template.Must(template.ParseFiles(dir + "/base.html")),
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
	if funcs != nil {
		for k, v := range funcs {
			ts.funcs[k] = v
		}
	}
	ts.base.Funcs(ts.funcs)
	return ts
}

func (ts *TemplateStore) Cache(name ...string) {
	ts.Lock()
	for i := 0; i < len(name); i++ {
		fmt.Printf("Cacheing %s\n", name[i])
		//ts.cache[name[i]] = template.Must(ts.base.ParseFiles(ts.dir+"/base.html", ts.dir+"/"+name[i]))
		ts.cache[name[i]] = template.Must(ts.base.ParseFiles(ts.dir + "/" + name[i]))
	}
	ts.Unlock()
}

func (ts *TemplateStore) Render(w http.ResponseWriter, name string, model interface{}) {
	t1 := time.Now().UnixNano()
	t, ok := ts.cache[name]
	if !ok {
		t = template.Must(ts.base.ParseFiles(ts.dir + "/" + name))
		ts.cache[name] = t
		t.Execute(w, model)
		fmt.Printf("Took %d ns\n", time.Now().UnixNano()-t1)
		fmt.Println("1")
		return
	}
	if changed(ts.dir + "/" + name) {
		t = template.Must(template.ParseFiles(ts.dir+"/base.html", ts.dir+"/"+name))
		ts.cache[name] = t
		t.Execute(w, model)
		fmt.Printf("Took %d ns\n", time.Now().UnixNano()-t1)
		fmt.Println("2")
		return
	}
	t.Execute(w, model)
	fmt.Printf("Took %d ns\n", time.Now().UnixNano()-t1)
	fmt.Println("3")
}

func changed(path string) bool {
	// gather file status
	fstat, err := os.Stat(path)
	// err check
	if err != nil {
		// if err; print timestamp and err, then panic
		log.Panic(err)
		return false
	}
	// no err; eval and return (true) file been modified within the last n seconds
	fmt.Printf("file time: %d, now - 3s %d\n", fstat.ModTime().Unix(), time.Now().Unix()-3)
	return fstat.ModTime().Unix() >= time.Now().Unix()-3
}

// func RenderNoBase(tmpl string) {
// 	t = template.Must(template.New("tmpl").Funcs(ts.funcs).Parse(tmpl), err error)
// }

func ContentType(w http.ResponseWriter, typ string) {
	w.Header().Set("Content Type", typ)
}
