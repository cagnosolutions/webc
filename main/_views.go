package main

import (
    "github.com/cagnosolutions/web/tmpl"
)

// default
var ts = *tmpl.Templates{}

var ts = *tmpl.Templates()


// custom config
var ts = *tmpl.Templates{
    Development: true,
    // optionally setup development enviroment
}




var tmplsMap = make(map[string]*template.Template)

type Templates struct {
    Development bool
    tmpls map[string]*template.Template
    funcs template.FuncMap
}
