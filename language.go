package gocloc

import (
	"path/filepath"
)

// Language is a type used to definitions and store statistics for one programming language.
type Language struct {
	Name         string
	lineComments []string
	Files        []string
	Code         int32
	Comments     int32
	Blanks       int32
	Total        int32
}

// Exts is the definition of the language name, keyed by the extension for each language.
var Exts = map[string]string{
	"go": "Go",
}

func getFileType(path string, opts *ClocOptions) (ext string, ok bool) {
	ext = filepath.Ext(path)

	if len(ext) >= 2 {
		return ext[1:], true
	}
	return ext, ok
}

// NewLanguage create language data store.
func NewLanguage(name string, lineComments []string) *Language {
	return &Language{
		Name:         name,
		lineComments: lineComments,
		Files:        []string{},
	}
}

// DefinedLanguages is the type information for mapping language name(key) and NewLanguage.
type DefinedLanguages struct {
	Langs map[string]*Language
}

// NewDefinedLanguages create DefinedLanguages.
func NewDefinedLanguages() *DefinedLanguages {
	return &DefinedLanguages{
		Langs: map[string]*Language{
			"Go": NewLanguage("Go", []string{"//"}),
		},
	}
}
