package main

import (
	_ "embed"
	"text/template"
)

var tmpl *template.Template

//go:embed templates/class
var classTemplate string

//go:embed templates/property
var propertyTemplate string

//go:embed templates/generic
var genericTemplate string

//go:embed templates/aacircular
var a_aTemplate string

//go:embed templates/initLinq
var initlinq string

func loadTemplates() (err error) {
	if tmpl == nil {
		tmpl, err = template.New("class").Parse(classTemplate)
		if err != nil {
			return
		}

		_, err = tmpl.New("property").Parse(propertyTemplate)
		if err != nil {
			return
		}

		_, err = tmpl.New("generic").Parse(genericTemplate)
		if err != nil {
			return
		}

		_, err = tmpl.New("aacircular").Parse(a_aTemplate)
		if err != nil {
			return
		}

		_, err = tmpl.New("initlinq").Parse(initlinq)
		if err != nil {
			return
		}
	}
	return
}
