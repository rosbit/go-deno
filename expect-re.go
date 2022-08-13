package deno

import (
	"regexp"
)

var (
	ignoreRE   = regexp.MustCompile(`^Deno \d+\.\d+\.\d+[\r\n]`)
	ignore2RE  = regexp.MustCompile(`^exit using .*?[\r\n]`)
	promptRE   = regexp.MustCompile(`^\x1b\x5b\x3f.*?\r\x1b\x5b.*?> \r\x1b\x5b.*`)  // "> "      prompt
	errLoadRE  = regexp.MustCompile(`^error in \-\-eval\-file file .*?[\r\n]`)     // "--eval-file error"  error message
	errRE      = regexp.MustCompile(`^Uncaught .+?Error: .+?[\r\n]`)   // error message
	errAtRE    = regexp.MustCompile(`^\s+at .+?[\r\n]`)                 // "  at xxxx"   error at xxxx
	blankRE    = regexp.MustCompile(`^[\r\n]`)
	goalRE     = regexp.MustCompile(`^\r?\x1b\x5b.*?[\r\n]?`)  // goal(xxx,xxx)
	nonResRE   = regexp.MustCompile(`\x1b\x5b.*?\n`)
	resultRE   = regexp.MustCompile(`^.+?[\r\n]`)

	jsonNameRE = regexp.MustCompile(`([\{\,] )([^ ]*?): `) // JSON name without quote
	jsonNameRepl = []byte(`${1}"${2}":`)
	functionId = []byte("[Function: ")
	nan = []byte("NaN")
	undefined = []byte("undefined")
)

