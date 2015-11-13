package webc

import (
    "fmt"
    "html/template"
    "net/http"
    "path/filepath"
    "log"
    "bytes"
)

var templates map[string]*template.Template
var bufpool *BufferPool

// load templates on program initialization
func init() {
    if templates == nil {
        templates = make(map[string]*template.Template)
    }

    templatesDir := "templates/"

    layouts, err := filepath.Glob(templatesDir + "layouts/*.tmpl")
    if err != nil {
        log.Fatal(err)
    }

    includes, err := filepath.Glob(templatesDir + "includes/*.tmpl")
    if err != nil {
        log.Fatal(err)
    }

    // generate our template map from layouts/* && includes/* dirs
    for _, layout := range layouts {
        files := append(includes, layout)
        templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
    }

    bufpool = BufferPoolInstance(64)

}

// render template is a wrapper around template.ExecuteTemplate()
// write into byte.buffer before writing to the responseWriter to catch any errors resulting from populating the template
func RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
    // ensure template exists in templates map
    tmpl, ok := templates[name]
    if !ok {
        return fmt.Errorf("The template %s does not exist.", name)
    }

    // create a buffer to temporarily write to and check for errors
    buf := bufpool.Get()
    defer bufpool.Put(buf)

    err := tmpl.ExecuteTemplate(buf, "base", data)
    if err != nil {
        return err
    }

    // ser header and write buffer to the responseWriter
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    buf.WriteTo(w)
    return nil
}

type BufferPool struct {
    ch chan *bytes.Buffer
}

func BufferPoolInstance(size int) *BufferPool {
    return &BufferPool{
        ch: make(chan *bytes.Buffer, size),
    }
}

func (bp *BufferPool) Get() *bytes.Buffer {
    var b *bytes.Buffer
    select {
        case b = <-bp.ch:
        default:
            b = bytes.NewBuffer([]byte{})
    }
    return b
}

func (bp *BufferPool) Put(b *bytes.Buffer) {
    b.Reset()
    select {
        case bp.ch <-b:
        default:
            // discard buffer if pool is full
    }
}
