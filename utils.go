package gocloc

import (
	"os"
	"path/filepath"
	"strings"
)

func containsComment(line string, multiLines [][]string) bool {
	for _, mlcomm := range multiLines {
		for _, comm := range mlcomm {
			if strings.Contains(line, comm) {
				return true
			}
		}
	}
	return false
}

func nextRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}

func checkOptionMatch(path string, info os.FileInfo, opts *ClocOptions) bool {
	// check match directory & file options
	dir := filepath.Dir(path)
	if opts.ReMatchDir != nil && !opts.ReMatchDir.MatchString(dir) {
		return false
	}

	return true
}

// getAllFiles return all of the files to be analyzed in paths.
func getAllFiles(paths []string, languages *DefinedLanguages, opts *ClocOptions) (result map[string]*Language, err error) {
	result = make(map[string]*Language, 0)

	for _, root := range paths {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// check match & not-match directory
			if match := checkOptionMatch(path, info, opts); !match {
				return nil
			}

			if ext, ok := getFileType(path, opts); ok {
				if targetExt, ok := Exts[ext]; ok {
					// check exclude extension
					if _, ok := result[targetExt]; !ok {
						result[targetExt] = NewLanguage(
							languages.Langs[targetExt].Name,
							languages.Langs[targetExt].lineComments)
					}
					result[targetExt].Files = append(result[targetExt].Files, path)
				}
			}
			return nil
		})
	}
	return
}
