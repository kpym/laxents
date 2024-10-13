package main

import (
	"bufio"
	"io"

	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

// isAscii returns false if there is a non-ASCII character in the data.
func isAscii(data []byte) bool {
	for _, b := range data {
		if b > 127 || b == 0 {
			return false
		}
	}
	return true
}

// guessUTF16 returns the "UTF-16 LE", "UTF-16 BE" if it looks like a valid UTF-16 or UTF-8.
// - first it checks if the data starts with a BOM
// - if no bom is found it counts the number of
//   - <null><ascii> pairs (for UTF-16 BE)
//   - <ascii><null> pairs (for UTF-16 LE)
func guessUTF16(data []byte) string {
	if len(data) >= 2 {
		switch {
		case data[0] == 0xFE && data[1] == 0xFF:
			return "UTF-16 BE"
		case data[0] == 0xFF && data[1] == 0xFE:
			return "UTF-16 LE"
		}
	}
	utf16be := 0
	for i := 0; i < len(data)-1; i += 2 {
		if data[i] == 0 && data[i+1] < 128 {
			utf16be++
		}
	}
	utf16le := 0
	for i := 0; i < len(data)-1; i += 2 {
		if data[i] < 128 && data[i+1] == 0 {
			utf16le++
		}
	}
	if utf16be > 0 || utf16le > 0 {
		if utf16be > utf16le {
			return "utf-16be"
		} else {
			return "utf-16le"
		}
	}
	return ""
}

// detectCharset returns the encoding of the data.
// if the data is Ascii it returns an empty string.
func detectCharset(data []byte) (encoding string) {
	if isAscii(data) {
		return ""
	}
	encoding = guessUTF16(data)
	if encoding != "" {
		return encoding
	}
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(data)
	if err != nil {
		return ""
	}
	return result.Charset
}

// utfReader returns a reader that converts the input to UTF-8
// if it is not already encoded in UTF-8.
// If an error occurs it returns the original reader.
func utfReader(r io.Reader) io.Reader {
	// New bufferedReader
	br := bufio.NewReader(r)

	// Peak the first 4096 bytes (at most)
	beginning, _ := br.Peek(4 * 1024)
	encoding := detectCharset(beginning)
	if encoding == "" || encoding == "UTF-8" {
		return br
	}
	// Convert the reader to UTF-8
	utfReader, err := charset.NewReaderLabel(encoding, br)
	if err != nil {
		return br
	}
	return utfReader
}
