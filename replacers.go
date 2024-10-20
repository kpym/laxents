package main

import (
	"bufio"
	"io"
	"strings"
)

// CopyWithReplace reads from an io.Reader, applies replacements, and writes to an io.Writer
func CopyWithReplace(r io.Reader, w io.Writer, replacer *strings.Replacer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		replacedLine := replacer.Replace(line)
		_, err := w.Write([]byte(replacedLine + "\n"))
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func toUnicodeLetters() *strings.Replacer {
	return strings.NewReplacer(
		"\\L", "Ł",
		"\\l", "ł",
		"\\O", "Ø",
		"\\o", "ø",
		"\\i", "ı",
		"\\j", "ȷ",
		"\\AA", "Å",
		"\\aa", "å",
		"\\AE", "Æ",
		"\\ae", "æ",
		"\\OE", "Œ",
		"\\oe", "œ",
		"\\ss", "ß",
	)
}

func toLatexLetters() *strings.Replacer {
	return strings.NewReplacer(
		"Ł", "{\\L}",
		"ł", "{\\l}",
		"Ø", "{\\O}",
		"ø", "{\\o}",
		"ı", "{\\i}",
		"ȷ", "{\\j}",
		"Å", "{\\AA}",
		"å", "{\\aa}",
		"Æ", "{\\AE}",
		"æ", "{\\ae}",
		"Œ", "{\\OE}",
		"œ", "{\\oe}",
		"ß", "{\\ss}",
	)
}
