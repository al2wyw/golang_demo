package main

import (
	"html/template"
	"log"
	"os"
	"testing"
)

var md = `
people name(out scope): {{ .Name }}
dog name(out scope): {{ .MyDog.Name }}
{{- with .MyDog }}
dog name(in scope): {{ .Name }} 
people name(in scope): {{ $.Name }}
{{ end }}
`

type People struct {
	Name  string
	Age   int
	MyDog Dog
}

type Dog struct {
	Name string
}

func TestTemplate(t *testing.T) {
	tpl := template.Must(template.New("first").Parse(md))
	p := People{Name: "Lucy", MyDog: Dog{Name: "Tom"}}
	if err := tpl.Execute(os.Stdout, p); err != nil {
		log.Fatal(err)
	}
}
