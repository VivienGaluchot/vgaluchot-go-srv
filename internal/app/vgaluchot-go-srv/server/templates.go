package server

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates map[string]map[string]*template.Template

func templateDir() string {
	return filepath.Join(os.Getenv("WEB_DIR"), "template")
}

func init() {
	var funcMap template.FuncMap = template.FuncMap{
		"title": strings.Title,
	}

	localizedTemplates := map[string]map[string][]string{
		"index.html": {
			"en": []string{"base.html", "index.html"},
			"fr": []string{"base.html", "base-fr.html", "index.html", "index-fr.html"},
		},
		"contact.html": {
			"en": []string{"base.html", "contact.html"},
			"fr": []string{"base.html", "base-fr.html", "contact.html", "contact-fr.html"},
		},
	}

	templates = make(map[string]map[string]*template.Template)

	for templateName, localizedSources := range localizedTemplates {
		if _, ok := templates[templateName]; ok != true {
			templates[templateName] = make(map[string]*template.Template)
		}
		for lang, sources := range localizedSources {
			sourcePath := make([]string, len(sources))
			for i, src := range sources {
				sourcePath[i] = filepath.Join(templateDir(), src)
			}
			templates[templateName][lang] = template.Must(
				template.New("index").Funcs(funcMap).ParseFiles(sourcePath...),
			)
		}
	}
}

// TemplateProcessingCtx is the input data to compute a template context
type TemplateProcessingCtx struct {
	templateName string
	request      *http.Request
}

func makeTemplateHandler(templateName string, getContext func(TemplateProcessingCtx) interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := getLanguageBase(r)

		template, ok := templates[templateName][lang]
		if !ok {
			log.Fatalf("Template %s for language %s not defined.", templateName, lang)
		}

		context := getContext(TemplateProcessingCtx{templateName: templateName, request: r})

		buf := &bytes.Buffer{}
		if err := template.ExecuteTemplate(buf, "base", context); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		buf.WriteTo(w)
	}
}
