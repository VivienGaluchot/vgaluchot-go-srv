package server

import (
	"net/http"

	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/conf"
)

var indexView func(w http.ResponseWriter, r *http.Request) = makeTemplateHandler("index.html", getIndexContext)
var contactView func(w http.ResponseWriter, r *http.Request) = makeTemplateHandler("contact.html", getContactContext)

// contexts

func getBaseContext(processingCtx TemplateProcessingCtx) interface{} {
	type ClientContext struct {
		Lang string
	}

	type BaseContext struct {
		Static       string
		Version      string
		RequestTime  string
		Client       ClientContext
		TemplateName string
	}

	lang := getLanguageBase(processingCtx.request)

	data := BaseContext{
		Static:  "/static",
		Version: conf.Version,
		Client: ClientContext{
			Lang: lang,
		},
		TemplateName: processingCtx.templateName,
	}
	return data
}

func getIndexContext(processingCtx TemplateProcessingCtx) interface{} {
	type Context struct {
		Base interface{}
	}
	data := Context{
		Base: getBaseContext(processingCtx),
	}
	return data
}

func getContactContext(processingCtx TemplateProcessingCtx) interface{} {
	type Context struct {
		Base interface{}
	}
	data := Context{
		Base: getBaseContext(processingCtx),
	}
	return data
}
