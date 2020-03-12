package reply

import (
	re2 "github.com/dlclark/regexp2"
)

const (
	delimiter = "d"
	embedded  = "b"
	empty     = "e"
	header    = "h"
	quote     = "q"
	signature = "s"
	text      = "t"
)

func classifyLine(line string) string {
	if isEmptyLine(line) {
		return empty
	}

	if isDelimiter(line) {
		return delimiter
	}

	if isSignature(line) {
		return signature
	}

	if isEmbeddedEmail(line) {
		return embedded
	}

	if isHeader(line) {
		return header
	}

	if isQuote(line) {
		return quote
	}

	return text
}

func isEmptyLine(line string) bool {
	ok, _ := re2.MustCompile(`^[[:blank:]]*$`, re2.RE2).MatchString(line)
	return ok
}

func isQuote(line string) bool {
	ok, _ := re2.MustCompile(`^[[:blank:]]*>`, re2.RE2).MatchString(line)
	return ok
}

func isDelimiter(line string) bool {
	ok, _ := re2.MustCompile(`^[[:blank:]]*[\-_,=+~#*ᐧ—]+[[:blank:]]*$`, re2.RE2).MatchString(line)
	return ok
}

func isSignature(line string) bool {
	// remove any markdown links
	stripped, _ := re2.MustCompile(`\[([^\]]+)\]\([^\)]+\)`, re2.RE2).Replace(line, "$1", 0, -1)
	for _, r := range patterns["SIGNATURE_REGEXES"] {
		ok, _ := r.MatchString(stripped)
		if ok {
			return true
		}
	}

	return false
}

func isHeader(line string) bool {
	for _, r := range patterns["EMAIL_HEADER_REGEXES"] {
		ok, _ := r.MatchString(line)
		if ok {
			return true
		}
	}

	return false
}

func isEmbeddedEmail(line string) bool {
	for _, r := range patterns["EMBEDDED_REGEXES"] {
		ok, _ := r.MatchString(line)
		if ok {
			return true
		}
	}

	return false
}
