package reply

import "strings"

// equivalent of "".slice!
func stringSliceBang(s string, i int) string {
	return strings.Join(sliceSliceBang(strings.Split(s, ""), i), "")
}

// equivalent of [].slice!
func sliceSliceBang(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

// equivalent of "abc"[0] = "x"
func stringReplaceChar(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

// equivalent of Ruby ""[start..end]
// .. is inclusive
// ... is exclusive
func sliceString(text string, start int, end int) string {
	var builder strings.Builder

	for i := start; i <= end; i++ {
		builder.WriteString(string(text[i]))
	}
	return builder.String()
}

// equivalent of Ruby [][start..end]
// .. is inclusive
// ... is exclusive
func sliceArray(lines []string, start int, end int) []string {
	newLines := []string{}
	for i := start; i <= end; i++ {
		newLines = append(newLines, lines[i])
	}
	return newLines
}
