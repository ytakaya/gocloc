package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	flags "github.com/jessevdk/go-flags"
	"github.com/ytakaya/gocloc"
)

// OutputTypeDefault is cloc's text output format for --output-type option
const OutputTypeDefault string = "default"

const fileHeader string = "File"
const languageHeader string = "Language"
const commonHeader string = "files          blank        comment           code"
const defaultOutputSeparator string = "-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------"

var rowLen = 79

// CmdOptions is gocloc command options.
// It is necessary to use notation that follows go-flags.
type CmdOptions struct {
	MatchDir string `long:"match-d" description:"include dir name (regex)"`
}

type outputBuilder struct {
	opts   *CmdOptions
	result *gocloc.Result
}

func newOutputBuilder(result *gocloc.Result, opts *CmdOptions) *outputBuilder {
	return &outputBuilder{
		opts,
		result,
	}
}

func (o *outputBuilder) WriteHeader() {
	headerLen := 28
	header := languageHeader
	fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
	fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, commonHeader)
	fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
}

func (o *outputBuilder) WriteFooter() {
	total := o.result.Total
	fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
	fmt.Printf("%-27v %6v %14v %14v %14v\n",
		"TOTAL", total.Total, total.Blanks, total.Comments, total.Code)
	fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
}

func (o *outputBuilder) WriteResult() {
	// write header
	o.WriteHeader()

	clocLangs := o.result.Languages

	var sortedLanguages gocloc.Languages
	for _, language := range clocLangs {
		if len(language.Files) != 0 {
			sortedLanguages = append(sortedLanguages, *language)
		}
	}
	sort.Sort(sortedLanguages)

	for _, language := range sortedLanguages {
		fmt.Printf("%-27v %6v %14v %14v %14v\n",
			language.Name, len(language.Files), language.Blanks, language.Comments, language.Code)
	}

	// write footer
	o.WriteFooter()
}

func main() {
	var opts CmdOptions
	clocOpts := gocloc.NewClocOptions()
	// parse command line options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gocloc"
	parser.Usage = "[OPTIONS] PATH[...]"

	paths, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	// value for language result
	languages := gocloc.NewDefinedLanguages()

	if len(paths) <= 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	if opts.MatchDir != "" {
		clocOpts.ReMatchDir = regexp.MustCompile(opts.MatchDir)
	}

	processor := gocloc.NewProcessor(languages, clocOpts)
	result, err := processor.Analyze(paths)
	if err != nil {
		fmt.Printf("fail gocloc analyze. error: %v\n", err)
		return
	}

	builder := newOutputBuilder(result, &opts)
	builder.WriteResult()
}
