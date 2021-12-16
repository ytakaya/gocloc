package gocloc

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// ClocLanguage is provide for xml-cloc and json format.
type ClocLanguage struct {
	Name       string `xml:"name,attr" json:"name,omitempty"`
	FilesCount int32  `xml:"files_count,attr" json:"files"`
	Code       int32  `xml:"code,attr" json:"code"`
	Comments   int32  `xml:"comment,attr" json:"comment"`
	Blanks     int32  `xml:"blank,attr" json:"blank"`
}

// Language is a type used to definitions and store statistics for one programming language.
type Language struct {
	Name         string
	lineComments []string
	multiLines   [][]string
	Files        []string
	Code         int32
	Comments     int32
	Blanks       int32
	Total        int32
}

// Languages is an array representation of Language.
type Languages []Language

func (ls Languages) Len() int {
	return len(ls)
}
func (ls Languages) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls Languages) Less(i, j int) bool {
	if ls[i].Code == ls[j].Code {
		return ls[i].Name < ls[j].Name
	}
	return ls[i].Code > ls[j].Code
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
func NewLanguage(name string, lineComments []string, multiLines [][]string) *Language {
	return &Language{
		Name:         name,
		lineComments: lineComments,
		multiLines:   multiLines,
		Files:        []string{},
	}
}

func lang2exts(lang string) (exts string) {
	es := []string{}
	for ext, l := range Exts {
		if lang == l {
			es = append(es, ext)
		}
	}
	return strings.Join(es, ", ")
}

// DefinedLanguages is the type information for mapping language name(key) and NewLanguage.
type DefinedLanguages struct {
	Langs map[string]*Language
}

// GetFormattedString return DefinedLanguages as a human readable string.
func (langs *DefinedLanguages) GetFormattedString() string {
	var buf bytes.Buffer
	printLangs := []string{}
	for _, lang := range langs.Langs {
		printLangs = append(printLangs, lang.Name)
	}
	sort.Strings(printLangs)
	for _, lang := range printLangs {
		buf.WriteString(fmt.Sprintf("%-30v (%s)\n", lang, lang2exts(lang)))
	}
	return buf.String()
}

// NewDefinedLanguages create DefinedLanguages.
func NewDefinedLanguages() *DefinedLanguages {
	return &DefinedLanguages{
		Langs: map[string]*Language{
			"Go": NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}}),
		},
	}
}
