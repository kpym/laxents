package transformers

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// toLaTeXAccents is a transformer that converts Unicode diacritics to LaTeX accents
type toLaTeXAccents struct {
	letter  rune
	accents []rune
}

// ToLaTeXAccents returns a transformer that converts Unicode diacritics to LaTeX accents
func ToLaTeXAccents() transform.Transformer {
	return &toLaTeXAccents{}
}

// Reset resets the transformer
func (t *toLaTeXAccents) Reset() {
	t.letter = 0
	t.accents = t.accents[:0]
}

// unicodeAccentsToLaTeX is a unicode to LaTeX accent mapping
// every diacritic is mapped to its LaTeX accent
var unicodeAccentsToLaTeX = map[rune]rune{
	0x300: '`',  // grave : à
	0x301: '\'', // acute : á
	0x302: '^',  // circumflex : â
	0x303: '~',  // tilde : ã
	0x304: '=',  // macron : ā
	0x306: 'u',  // breve : ă
	0x307: '.',  // dot-over : ġ
	0x308: '"',  // two dots (umlaut, diaeresis) : ä
	0x30A: 'r',  // ring : å
	0x30B: 'H',  // double acute (long Hungarian umlaut) : ő
	0x30C: 'v',  // háček : č
	0x323: 'd',  // dot-under : ẹ
	0x327: 'c',  // cedilla : ç
	0x328: 'k',  // ogonek : ą
	0x331: 'b',  // bar-under (macron below) : ḵ
}

// unicodeLettersToLaTeX is a unicode to LaTeX letter mapping
// every letter is mapped to its LaTeX equivalent
var unicodeLettersToLaTeX = map[rune]string{
	'Ł': "{\\L}",
	'ł': "{\\l}",
	'Ø': "{\\O}",
	'ø': "{\\o}",
	'ı': "{\\i}",
	'ȷ': "{\\j}",
	'Å': "{\\AA}",
	'å': "{\\aa}",
	'Æ': "{\\AE}",
	'æ': "{\\ae}",
	'Œ': "{\\OE}",
	'œ': "{\\oe}",
	'ß': "{\\ss}",
}

type adjusment struct {
	fromAccent rune
	fromLetter rune
	toAccent   rune
	toLetter   rune
}

var adjustments = []adjusment{
	{'`', 'i', '`', 'ı'},
	{'\'', 'i', '\'', 'ı'},
	{'^', 'i', '^', 'ı'},
	{'"', 'i', '"', 'ı'},
	{'`', 'j', '`', 'ȷ'},
	{'\'', 'j', '\'', 'ȷ'},
	{'^', 'j', '^', 'ȷ'},
	{'"', 'j', '"', 'ȷ'},
	{'"', 'j', '"', 'ȷ'},
	{'r', 'a', 0, 'å'},
	{'r', 'A', 0, 'Å'},
}

var adjustLetters = "ijaA"

// adjust adjusts the accents and letter according to the adjustments table
func (t *toLaTeXAccents) adjust() {
	if len(t.accents) == 0 || strings.IndexRune(adjustLetters, t.letter) < 0 {
		return
	}
	for _, a := range adjustments {
		if a.fromAccent == t.accents[0] && a.fromLetter == t.letter {
			if a.toAccent == 0 {
				t.accents = t.accents[1:]
			} else {
				t.accents[0] = a.toAccent
			}
			t.letter = a.toLetter
			return
		}
	}
}

// writeLaTeXRune writes r to dst and increments nDst by the number of bytes written
// it returns true if the whole rune was written
func writeLaTeXRune(dst []byte, r rune, nDst *int) bool {
	if s, ok := unicodeLettersToLaTeX[r]; ok {
		return write(dst, s, nDst)
	}
	return writeRune(dst, r, nDst)
}

// writeLaTeXLetter writes the LaTeX letter for a rune to dst
// it returns true if the letter was written or if there is no letter to write (t.letter == 0)
// if it returns false, nDst is not modified
func (t *toLaTeXAccents) writeLaTeXLetter(dst []byte, nDst *int, inGroup bool) (done bool) {
	if t.letter == 0 {
		return true
	}
	// reset the letter if it was written
	defer func() {
		if done {
			t.letter = 0
		}
	}()
	if s, ok := unicodeLettersToLaTeX[t.letter]; ok {
		return write(dst, s, nDst)
	}
	if inGroup {
		return write(dst, fmt.Sprintf("{%c}", t.letter), nDst)
	}
	return writeLaTeXRune(dst, t.letter, nDst)
}

// writeLaTeXAccent writes the commulated LaTeX accents followed by the letter to dst
// it returns true if everything was written
// if it returns false, nDst is not modified
func (t *toLaTeXAccents) writeLaTeXAccent(dst []byte, nDst *int) (done bool) {
	n := *nDst
	inGroup := false
	// adjust the accents and letter
	t.adjust()
	// write the accents
	for i := len(t.accents) - 1; i >= 0; i-- {
		if !writeRune(dst, '\\', &n) || !writeRune(dst, t.accents[i], &n) {
			return false
		}
		inGroup = isLatin(t.accents[i])
	}
	// write the letter (and reset it)
	if !t.writeLaTeXLetter(dst, &n, inGroup) {
		return false
	}
	// reset the accents
	t.accents = t.accents[:0]
	// everything was written
	*nDst = n
	return true
}

// Transform converts Unicode diacritics to LaTeX accents
// src is supposed to be a valid UTF-8 string in NFD form
func (t *toLaTeXAccents) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	// prev is the previous rune in src (of size psize)
	// it is checked to see if it is a diacritic for the current rune
	var (
		r    rune // current rune (temp variable)
		size int  // size of the current rune (in bytes, temp variable)
	)
	// loop over the runes in src
	for nSrc < len(src) {
		// read the next rune
		r, size = utf8.DecodeRune(src[nSrc:])
		if r == utf8.RuneError {
			err := transform.ErrShortSrc
			if atEOF {
				err = fmt.Errorf("invalid UTF-8 encoding")
			}
			return nDst, nSrc, err
		}
		// check if the rune is a diacritic
		if accent, ok := unicodeAccentsToLaTeX[r]; ok {
			t.accents = append(t.accents, accent)
		} else {
			// write commulated accents followed by the letter
			if !t.writeLaTeXAccent(dst, &nDst) {
				return nDst, nSrc, transform.ErrShortDst
			}
			// save the current rune as the letter for the next accents (if any)
			t.letter = r
		}
		nSrc += size
	}
	if len(t.accents) > 0 && !atEOF {
		// we are collecting accents, but we reached the end of src
		return nDst, nSrc, transform.ErrShortSrc
	}
	if !t.writeLaTeXAccent(dst, &nDst) {
		return nDst, nSrc, transform.ErrShortDst
	}
	return nDst, nSrc, nil
}
