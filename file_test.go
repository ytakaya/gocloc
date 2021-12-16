package gocloc

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestAnalayzeFile4Python(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.py")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3
	"""
	pass
`))

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 4 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4PythonInvalid(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.py")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3"""
	pass
`))

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 3 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4PythonNoShebang(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.py")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`a = '''hello
	world
	'''

	b = 1
	"""hello
	commen
	"""

	print a, b
`))

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 2 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 3 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 5 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4Go(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.go")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`package main

func main() {
	var n string /*
		comment
		comment
	*/
}
`))

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 3 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 4 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4GoWithOnelineBlockComment(t *testing.T) {
	t.SkipNow()
	tmpfile, err := ioutil.TempFile("", "tmp.go")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`package main

func main() {
	st := "/*"
	a := 1
	en := "*/"
	/* comment */
}
`))

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 1 { // cloc->3, tokei->1, gocloc->4
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 6 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4GoWithCommentInnerBlockComment(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.go")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`package main

func main() {
	// comment /*
	a := 1
	b := 2
}
`))

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 1 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 5 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4GoWithNoComment(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.go")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`package main

	func main() {
		a := "/*                */"
		b := "//                  "
	}
`))

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 0 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 5 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4ATSWithDoubleMultilineComments(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.java")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`/* com */
(* co *)

vo (*
com *)

vo /*
jife */

vo /* ff */
`))

	language := NewLanguage("ATS", []string{"//"}, [][]string{{"(*", "*)"}, {"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 3 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 4 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "ATS" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4JavaWithCommentInCodeLine(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.java")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`public class Sample {
		public static void main(String args[]){
		int a; /* A takes care of counts */
		int b;
		int c;
		String d; /*Just adding comments */
		bool e; /*
		comment*/
		bool f; /*
		comment1
		comment2
		*/
		/*End of Main*/
		}
		}
`))

	language := NewLanguage("Java", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 0 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 5 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 10 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Java" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4Makefile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "Makefile.am")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`# This is a simple Makefile with comments
	.PHONY: test build

	build:
		mkdir -p bin
		GO111MODULE=on go build -o ./bin/gocloc cmd/gocloc/main.go

	# Another comment
	update-package:
		GO111MODULE=on go get -u github.com/ytakaya/gocloc

	run-example:
		GO111MODULE=on go run examples/languages.go
		GO111MODULE=on go run examples/files.go

	test:
		GO111MODULE=on go test -v
`))

	language := NewLanguage("Makefile", []string{"#"}, [][]string{{"", ""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 4 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 2 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 11 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Makefile" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3
	"""
	pass
`))

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeReader("test.py", language, buf, clocOpts)

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 4 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}
