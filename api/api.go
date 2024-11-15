// api package is a package that contains the functions that
// convert between LaTeX accents and Unicode characters.
package api

import (
	"io"

	"github.com/kpym/laxents/transformers"
	"github.com/kpym/laxents/utfreader"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ToUnicode converts LaTeX accents to Unicode characters.
func ToUnicode(out io.Writer, in io.Reader) error {
	in = utfreader.New(in, utfreader.NFC)
	outt := transform.NewWriter(out, transform.Chain(transformers.ToUnicodeAccents(), norm.NFC))
	defer outt.Close()

	_, err := io.Copy(outt, in)
	return err
}

// ToLaTeX converts Unicode characters to LaTeX accents.
func ToLaTeX(out io.Writer, in io.Reader) error {
	in = utfreader.New(in, utfreader.NFD)
	outt := transform.NewWriter(out, transform.Chain(transformers.ToLaTeXAccents(), norm.NFC))
	defer outt.Close()

	_, err := io.Copy(outt, in)
	return err
}
