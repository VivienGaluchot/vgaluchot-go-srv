package server

import (
	"log"
	"net/http"

	"golang.org/x/text/language"
)

// Supported languages
var langMatcher = language.NewMatcher([]language.Tag{
	language.English,
	language.French,
})

func getLanguageTag(r *http.Request) language.Tag {
	lang, _ := r.Cookie("lang")
	accept := r.Header.Get("Accept-Language")
	var tag language.Tag
	if lang != nil {
		tag, _ = language.MatchStrings(langMatcher, lang.Value, accept)
	} else {
		tag, _ = language.MatchStrings(langMatcher, accept)
	}
	return tag
}

func getLanguageBase(r *http.Request) string {
	tag := getLanguageTag(r)
	base, _ := tag.Base()
	str := base.String()
	if str != "en" && str != "fr" {
		log.Printf("Unexpected language base '%s', fallback to 'en'", str)
		str = "en"
	}
	return str
}
