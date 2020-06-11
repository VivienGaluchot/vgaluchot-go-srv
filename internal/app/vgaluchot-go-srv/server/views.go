package server

import (
	"net/http"

	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/conf"
)

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

// index

func getIndexContext(processingCtx TemplateProcessingCtx) interface{} {
	type Context struct {
		Base interface{}
	}
	data := Context{
		Base: getBaseContext(processingCtx),
	}
	return data
}

var indexView func(w http.ResponseWriter, r *http.Request) = makeTemplateHandler("index.html", getIndexContext)

// portfolio

func getPortfolioContext(processingCtx TemplateProcessingCtx) interface{} {
	type Context struct {
		Base interface{}
	}
	data := Context{
		Base: getBaseContext(processingCtx),
	}
	return data
}

var portfolioView func(w http.ResponseWriter, r *http.Request) = makeTemplateHandler("portfolio.html", getPortfolioContext)

// contact

func getContactContext(processingCtx TemplateProcessingCtx) interface{} {
	type Context struct {
		Base interface{}
	}
	data := Context{
		Base: getBaseContext(processingCtx),
	}
	return data
}

var contactView func(w http.ResponseWriter, r *http.Request) = makeTemplateHandler("contact.html", getContactContext)
