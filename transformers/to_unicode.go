package transformers

import (
	"bytes"
	"strings"

	"golang.org/x/text/transform"
)

// toUnicodeAccents is a transformer that converts LaTeX accents to Unicode diacritics
type toUnicodeAccents struct {
	printBracket bool
	letter       rune
	accents      []rune
}

// ToLaTeXAccents returns a transformer that converts Unicode diacritics to LaTeX accents
func ToUnicodeAccents() transform.Transformer {
	return &toUnicodeAccents{}
}

// Reset resets the transformer
func (t *toUnicodeAccents) Reset() {
	t.printBracket = false
	t.letter = 0
	t.accents = t.accents[:0]
}

func (t *toUnicodeAccents) isZero() bool {
	return !t.printBracket && t.letter == 0 && len(t.accents) == 0
}

// startGroup sets the startGroup flag
func (t *toUnicodeAccents) startGroup(c byte) bool {
	if c == '{' && t.isZero() {
		t.printBracket = true
		return true
	}
	return false
}

// write writes the utf8 letter ans accents to dst.
func (t *toUnicodeAccents) write(dst []byte, nDst *int) (ok bool) {
	n := *nDst
	if t.printBracket {
		if !writeByte(dst, '{', &n) {
			return false
		}
	}
	if t.letter != 0 {
		if !writeRune(dst, t.letter, &n) {
			return false
		}
	}
	for i := len(t.accents) - 1; i >= 0; i-- {
		if !writeRune(dst, t.accents[i], &n) {
			return false
		}
	}
	// everything was written, reset the transformer
	*nDst = n
	t.Reset()
	return true
}

type latexSpecialType int

const (
	latexSpecialNone latexSpecialType = iota
	latexSpecialLetter
	latexSpecialNonLetterAccent
	latexSpecialLetterAccent
)

type latexSpecial struct {
	spType latexSpecialType
	utf8   rune
}

var latexToUnicode = map[string]latexSpecial{
	// Non-letter acccents
	"`":  {latexSpecialNonLetterAccent, 0x300}, // grave : à
	"'":  {latexSpecialNonLetterAccent, 0x301}, // acute : á
	"^":  {latexSpecialNonLetterAccent, 0x302}, // circumflex : â
	"~":  {latexSpecialNonLetterAccent, 0x303}, // tilde : ã
	"=":  {latexSpecialNonLetterAccent, 0x304}, // macron : ā
	"u":  {latexSpecialNonLetterAccent, 0x306}, // breve : ă
	".":  {latexSpecialNonLetterAccent, 0x307}, // dot-over : ġ
	"\"": {latexSpecialNonLetterAccent, 0x308}, // two dots (umlaut, diaeresis) : ä
	// Letter accents
	"r": {latexSpecialLetterAccent, 0x30A}, // ring : å
	"H": {latexSpecialLetterAccent, 0x30B}, // double acute (long Hungarian umlaut) : ő
	"v": {latexSpecialLetterAccent, 0x30C}, // háček : č
	"d": {latexSpecialLetterAccent, 0x323}, // dot-under : ẹ
	"c": {latexSpecialLetterAccent, 0x327}, // cedilla : ç
	"k": {latexSpecialLetterAccent, 0x328}, // ogonek : ą
	"b": {latexSpecialLetterAccent, 0x331}, // bar-under (macron below) : ḵ
	// Special letters
	"L":  {latexSpecialLetter, 'Ł'},
	"l":  {latexSpecialLetter, 'ł'},
	"O":  {latexSpecialLetter, 'Ø'},
	"o":  {latexSpecialLetter, 'ø'},
	"i":  {latexSpecialLetter, 'ı'},
	"j":  {latexSpecialLetter, 'ȷ'},
	"AA": {latexSpecialLetter, 'Å'},
	"aa": {latexSpecialLetter, 'å'},
	"AE": {latexSpecialLetter, 'Æ'},
	"ae": {latexSpecialLetter, 'æ'},
	"OE": {latexSpecialLetter, 'Œ'},
	"oe": {latexSpecialLetter, 'œ'},
	"ss": {latexSpecialLetter, 'ß'},
}

var noneLatexSpecial = latexSpecial{spType: latexSpecialNone}

const (
	nonletteraccent string = "`'^~=.u\""
	firstOfTwo      string = "Aes" // no need to put Oo because they are letter accents
)

// getSpecial start looking for a latex special at the beggining of the src.
// It returns true for needMore if src is empty or if the src is a single character
// that is the beginning of a two character special.
// In this case, it returns noneLatexSpecial and 0.
// It also returns true for needMore if a special is found at the end of src
// and it is not obvious it is not the beginning of some other macro.
// But in this case, it returns the special and the number of bytes read,
// because we can be at the end of the file.
// If needMore is false, it returns the special and the number of bytes read.
// If no special is found, it returns noneLatexSpecial and 0.
// If the special is not a non-letter accent, it gobbles the next space if it is there.
func getSpecial(src []byte) (ls latexSpecial, n int, needMore bool) {
	if len(src) == 0 {
		return noneLatexSpecial, 0, true
	}
	// if is a non-letter accent
	if strings.IndexByte(nonletteraccent, src[0]) >= 0 {
		return latexToUnicode[string(src[0])], 1, false
	}
	// get the longest possible latex macro name
	i := 0
	for ; i < len(src); i++ {
		if !isLatinByte(src[i]) && src[i] != '@' {
			break
		}
	}
	if ls, ok := latexToUnicode[string(src[:i])]; ok {
		if ls.spType == latexSpecialLetterAccent || ls.spType == latexSpecialLetter {
			if i < len(src) && src[i] == ' ' {
				// gobble the next space
				return ls, i + 1, false
			}
			if ls.spType == latexSpecialLetter && i < len(src) && src[i] == '{' {
				if i+1 == len(src) {
					// maybe we should gobble the empty group
					return ls, i, true
				}
				if src[i+1] == '}' {
					// gobble the empty group
					return ls, i + 2, false
				}
			}
			return ls, i, i == len(src)
		}
		return ls, i, false
	}
	if len(src) == 1 {
		// if the src is a single character that is the beginning of a two character special
		// return needMore = true
		return noneLatexSpecial, 0, strings.IndexByte(firstOfTwo, src[0]) >= 0
	}
	return noneLatexSpecial, 0, false
}

// getLetter check if the beginning of the src is a letter or {letter}.
// If it is a letter or {letter}, it returns the letter and the number of bytes read (1 or 3).
// If it is not a letter or {letter}, it returns 0,0, flase.
// If the src is "{" or "{letter", it returns 0,0, true.
func getLetter(src []byte) (l rune, n int, needMore bool) {
	if len(src) == 0 {
		return 0, 0, true
	}
	if src[0] == '{' {
		if len(src) == 1 {
			return 0, 0, true
		}
		if src[1] == '}' {
			return 0, 2, false
		}
		if len(src) == 2 {
			return 0, 0, true
		}
		if isLatinByte(src[1]) && src[2] == '}' {
			return rune(src[1]), 3, false
		}
		if src[1] == '\\' {
			ls, n, more := getSpecial(src[2:])
			if more {
				return 0, 0, true
			}
			if ls.spType != latexSpecialLetter || 2+n == len(src) || src[2+n] != '}' {
				return 0, 0, false
			}
			return ls.utf8, 3 + n, false
		}
		return 0, 0, false
	}
	if !isLatinByte(src[0]) {
		// we do not treat latexSpecialLetter here
		return 0, 0, false
	}
	return rune(src[0]), 1, false
}

// Transform converts LaTeX accents to Unicode diacritics
// src is supposed to be a valid UTF-8 string in NFD form
func (t *toUnicodeAccents) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		if src[nSrc] != '\\' {
			if t.startGroup(src[nSrc]) {
				nSrc++
				continue
			}
			if t.printBracket && src[nSrc] == '}' {
				t.printBracket = false
				nSrc++
			}
			// write collected accents to dst
			if !t.write(dst, &nDst) {
				// not enough space in dst
				return nDst, nSrc, transform.ErrShortDst
			}
			// find the next \ in src
			i := bytes.IndexByte(src[nSrc:], '\\')
			if i < 0 {
				i = len(src) - nSrc
			}
			if !writeBytes(dst, src[nSrc:nSrc+i], &nDst) {
				// not enough space in dst
				return nDst, nSrc, transform.ErrShortDst
			}
			nSrc += i
			continue
		}
		// get the special
		sp, n, needMore := getSpecial(src[nSrc+1:])
		if needMore && !atEOF {
			// we need more data to know how to process the special
			return nDst, nSrc, transform.ErrShortSrc
		}
		if sp.spType == latexSpecialNone {
			// n = 0 here
			// write the accents (without letter) to dst
			if !t.write(dst, &nDst) {
				// not enough space in dst
				return nDst, nSrc, transform.ErrShortDst
			}
			// can we get the escaped (after '\') character?
			if nSrc+1 >= len(src) && !atEOF {
				// we need more data to know how to process the special
				return nDst, nSrc, transform.ErrShortSrc
			}
			// write the `\` to dst
			if !writeByte(dst, '\\', &nDst) {
				// not enough space in dst
				return nDst, nSrc, transform.ErrShortDst
			}
			if !writeByte(dst, src[nSrc+1], &nDst) {
				// not enough space in dst
				return nDst, nSrc, transform.ErrShortDst
			}
			nSrc += 2
			continue
		}
		n++
		var m int
		if sp.spType == latexSpecialLetter {
			t.letter = sp.utf8
		} else {
			// get the letter
			t.letter, m, needMore = getLetter(src[nSrc+n:])
			if needMore && !atEOF {
				// we need more data to know how to process the letter
				return nDst, nSrc, transform.ErrShortSrc
			}
			t.accents = append(t.accents, sp.utf8)
			n += m
		}
		nSrc += n
		if t.printBracket {
			if nSrc >= len(src) && !atEOF {
				// we need more data to know how to process the letter
				return nDst, nSrc, transform.ErrShortSrc
			}
			if src[nSrc] == '}' {
				t.printBracket = false
				nSrc++
			}
		}
		if t.letter != 0 {
			if !t.write(dst, &nDst) {
				// not enough space in dst
				return nDst, nSrc, transform.ErrShortDst
			}
		}
	}
	if !t.write(dst, &nDst) {
		// not enough space in dst
		return nDst, nSrc, transform.ErrShortDst
	}
	return nDst, nSrc, nil
}
