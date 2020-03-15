package reply

import (
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

func init() {
	// The default configuration is set to 'forever'.
	// I am not expecting any regex to take more than a few milliseconds.
	// Setting this value to 1s just to be on the safe side.
	regexp2.DefaultMatchTimeout = 1 * time.Second
}

// FromReader returns the reply text from the e-mail text body.
func FromReader(reader io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return FromText(string(bytes)), nil
}

// FromText returns the reply text from the e-mail text body.
func FromText(text string) string {
	if strings.TrimSpace(text) == "" {
		return text
	}

	// do some cleanup
	text = cleanup(text)

	// from now on, we'll work on a line-by-line basis
	lines := strings.Split(text, "\n")
	patternBuilder := strings.Builder{}

	for _, line := range lines {
		patternBuilder.WriteString(classifyLine(line))
	}

	pattern := patternBuilder.String()

	// remove everything after the first delimiter
	{
		match, err := regexp2.MustCompile(`d`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = sliceString(pattern, 0, index-1)
			lines = sliceArray(lines, 0, index-1)
		}
	}

	// remove all mobile signatures
	for {
		match, err := regexp2.MustCompile(`s`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = stringSliceBang(pattern, index)
			lines = sliceSliceBang(lines, index)
		} else {
			break
		}
	}

	// when the reply is at the end of the email
	{
		match, err := regexp2.MustCompile(`^(b[^t]+)*b[bqeh]+t[et]*$`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			submatch, err := regexp2.MustCompile(`t[et]*$`, regexp2.RE2).FindStringMatch(pattern)
			if err != nil {
				return ""
			}

			index := submatch.Index
			pattern = ""
			lines = sliceArray(lines, index, len(lines)-1)
		}
	}

	// if there is an embedded email marker, not followed by a quote
	// then take everything up to that marker
	{
		match, err := regexp2.MustCompile(`te*b[^q]*$`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = sliceString(pattern, 0, index)
			lines = sliceArray(lines, 0, index)
		}
	}

	// if there is an embedded email marker, followed by a huge quote
	// then take everything up to that marker
	{
		match, err := regexp2.MustCompile(`te*b[eqbh]*([te]*)$`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil && strings.Count(match.GroupByNumber(1).String(), "t") < 7 {
			submatch, err := regexp2.MustCompile(`te*b[eqbh]*[te]*$`, regexp2.RE2).FindStringMatch(pattern)
			if err != nil {
				return ""
			}
			index := submatch.Index
			pattern = sliceString(pattern, 0, index)
			lines = sliceArray(lines, 0, index)
		}
	}

	// if there is some text before a huge quote ending the email,
	// then remove the quote
	{
		match, err := regexp2.MustCompile(`te*[qbe]+$`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = sliceString(pattern, 0, index)
			lines = sliceArray(lines, 0, index)
		}
	}

	// if there still are some embedded email markers, just remove them
	for {
		match, err := regexp2.MustCompile(`b`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = stringSliceBang(pattern, index)
			lines = sliceSliceBang(lines, index)
		} else {
			break
		}
	}

	// fix email headers when they span over multiple lines
	{
		match, err := regexp2.MustCompile(`h+[hte]+h+e`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			for i := 0; i < match.Length; i++ {
				c := []rune(header)[0]
				pattern = stringReplaceChar(pattern, c, index+i)
			}
		}
	}

	// if there are at least 3 consecutive email headers,
	// take everything up to these headers
	{
		match, err := regexp2.MustCompile(`t[eq]*h{3,}`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = sliceString(pattern, 0, index)
			lines = sliceArray(lines, 0, index)
		}
	}

	// if there still are some email headers, just remove them
	for {
		match, err := regexp2.MustCompile(`h`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match != nil {
			index := match.Index
			pattern = stringSliceBang(pattern, index)
			lines = sliceSliceBang(lines, index)
		} else {
			break
		}
	}

	// remove trailing quotes when there's at least one line of text
	{
		match1, err := regexp2.MustCompile(`t`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		match2, err := regexp2.MustCompile(`[eq]+$`, regexp2.RE2).FindStringMatch(pattern)
		if err != nil {
			return ""
		}
		if match1 != nil && match2 != nil {
			index := match2.Index
			pattern = sliceString(pattern, 0, index-1)
			lines = sliceArray(lines, 0, index-1)
		}
	}

	return strings.Join(lines, "\n")
}

func cleanup(text string) string {
	// normalize line endings
	replacer := strings.NewReplacer(
		"\r\n", "\n",
		"\r", "\n",
	)

	text = replacer.Replace(text)

	// remove PGP markers
	for _, r := range patterns["REMOVE_PGP_MARKERS_REGEX"] {
		text, _ = r.Replace(text, "", 0, -1)
	}

	// remove unsubscribe links
	for _, r := range patterns["REMOVE_UNSUBSCRIBE_REGEX"] {
		text, _ = r.Replace(text, "", 0, -1)
	}

	// remove alias-style quotes marker
	for _, r := range patterns["REMOVE_ALIAS_REGEX"] {
		text, _ = r.Replace(text, "", 0, -1)
	}

	// change enclosed-style quotes format
	for _, r := range patterns["CHANGE_ENCLOSED_QUOTE_ONE_REGEX"] {
		text, _ = r.ReplaceFunc(text, func(m regexp2.Match) string {
			newText, _ := regexp2.MustCompile(`^`, regexp2.RE2).Replace(m.GroupByNumber(2).String(), "> ", 0, -1)
			return newText
		}, 0, -1)
	}

	for _, r := range patterns["CHANGE_ENCLOSED_QUOTE_TWO_REGEX"] {
		text, _ = r.ReplaceFunc(text, func(m regexp2.Match) string {
			newText, _ := regexp2.MustCompile(`^`, regexp2.RE2).Replace(m.GroupByNumber(1).String(), "> ", 0, -1)
			return newText
		}, 0, -1)
	}

	// fix all quotes formats
	for _, r := range patterns["FIX_QUOTES_FORMAT_REGEX"] {
		text, _ = r.ReplaceFunc(text, func(m regexp2.Match) string {
			newText, _ := regexp2.MustCompile(`([[:alpha:]]+>|\|)`, regexp2.RE2).Replace(m.GroupByNumber(1).String(), ">", 0, -1)
			return newText
		}, 0, -1)
	}

	// fix embedded email markers that might span over multiple lines
	for _, regex := range patterns["FIX_EMBEDDED_REGEX"] {
		text, _ = regex.ReplaceFunc(text, func(m regexp2.Match) string {
			if strings.Count(m.String(), "\n") > 4 {
				return m.String()
			}
			newText, _ := regexp2.MustCompile(`\n+[[:space:]]*`, regexp2.RE2).Replace(m.String(), " ", 0, -1)
			return newText
		}, 0, -1)
	}

	// remove leading/trailing whitespaces
	return strings.TrimSpace(text)
}
