package controllers

import (
	log "github.com/sirupsen/logrus"
	"html/template"
)

func initTemplates() *template.Template {
	tmpl := template.New("")

	patterns := []string{
		"./views/blocks/*.html",
		"./views/pages/*.html",
	}

	for _, pattern := range patterns {
		globedTemplates, err := tmpl.ParseGlob(pattern)
		if err != nil {
			log.Fatalf("Failed to parse glob pattern '%s': %v", pattern, err)
		}
		tmpl = globedTemplates
	}

	return tmpl
}
