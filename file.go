package gocloc

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ClocFile is collecting to line count result.
type ClocFile struct {
	Code     int32  `xml:"code,attr" json:"code"`
	Comments int32  `xml:"comment,attr" json:"comment"`
	Blanks   int32  `xml:"blank,attr" json:"blank"`
	Name     string `xml:"name,attr" json:"name"`
	Lang     string `xml:"language,attr" json:"language"`
}

// ClocFiles is gocloc result set.
type ClocFiles []ClocFile

// AnalyzeFile is analyzing file, this function calls AnalyzeReader() inside.
func AnalyzeFile(filename string, language *Language, opts *ClocOptions) *ClocFile {
	fp, err := os.Open(filename)
	if err != nil {
		// ignore error
		return &ClocFile{Name: filename}
	}
	defer fp.Close()

	return AnalyzeReader(filename, language, fp, opts)
}

// AnalyzeReader is analyzing file for io.Reader.
func AnalyzeReader(filename string, language *Language, file io.Reader, opts *ClocOptions) *ClocFile {
	clocFile := &ClocFile{
		Name: filename,
		Lang: language.Name,
	}

	buf := getByteSlice()
	defer putByteSlice(buf)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lineOrg := scanner.Text()
		line := strings.TrimSpace(lineOrg)

		if len(strings.TrimSpace(line)) == 0 {
			clocFile.Blanks++
			continue
		}

		for _, singleComment := range language.lineComments {
			if strings.HasPrefix(line, singleComment) {
				clocFile.Comments++
				continue
			}
		}

		language.Code++
	}

	return clocFile
}
