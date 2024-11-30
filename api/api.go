// api package is a package that contains the functions that
// convert between LaTeX accents and Unicode characters.
package api

import (
	"io"

	"github.com/kpym/laxents/transformers"
	"golang.org/x/text/unicode/norm"

	"github.com/kpym/utf8reader"
)

// ToUnicode converts LaTeX accents to Unicode characters.
func ToUnicode(out io.Writer, in io.Reader) error {
	in = utf8reader.New(in,
		utf8reader.WithTransform(norm.NFC, transformers.ToUnicodeAccents(), norm.NFC))

	_, err := io.Copy(out, in)
	return err
}

// ToLaTeX converts Unicode characters to LaTeX accents.
func ToLaTeX(out io.Writer, in io.Reader) error {
	in = utf8reader.New(in,
		utf8reader.WithTransform(norm.NFD, transformers.ToLaTeXAccents(), norm.NFC))

	_, err := io.Copy(out, in)
	return err
}
