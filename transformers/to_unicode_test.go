package transformers

import (
	"bytes"
	"testing"
)

func TestToUnicodeAccents_Reset(t *testing.T) {
	lat := &toUnicodeAccents{letter: 'a', accents: []rune{'`'}}
	lat.Reset()
	if lat.letter != 0 {
		t.Errorf("expected letter=0, got letter=%v", lat.letter)
	}
	if len(lat.accents) != 0 {
		t.Errorf("expected accents=nil, got accents=%v", lat.accents)
	}
}

func TestToUnicodeAccents_write(t *testing.T) {
	data := []struct {
		t     *toUnicodeAccents // transformer
		dst   []byte            // destination slice
		ndst  int               // number of bytes already written to the destination slice
		exp   []byte            // expected destination slice after writing
		expn  int               // expected number of bytes written to the destination slice after writing
		expok bool              // expected return value of the write operation
	}{
		{&toUnicodeAccents{letter: 0, accents: []rune{}}, []byte("....."), 0, []byte("....."), 0, true},
		{&toUnicodeAccents{letter: 'a', accents: []rune{}}, []byte("....."), 0, []byte("a...."), 1, true},
		{&toUnicodeAccents{letter: 0, accents: []rune{0x300}}, []byte("....."), 0, []byte{0xCC, 0x80, '.', '.', '.'}, 2, true},
		{&toUnicodeAccents{letter: 0, accents: []rune{0x300, 0x301}}, []byte("....."), 0, []byte{0xCC, 0x81, 0xCC, 0x80, '.'}, 4, true},
		{&toUnicodeAccents{letter: 0, accents: []rune{0x300, 0x301, 0x302}}, []byte("....."), 0, []byte{0xCC, 0x82, 0xCC, 0x81, 0xCC}, 0, false},
		{&toUnicodeAccents{letter: 'a', accents: []rune{0x300}}, []byte("....."), 0, []byte{'a', 0xCC, 0x80, '.', '.'}, 3, true},
		{&toUnicodeAccents{letter: 'a', accents: []rune{0x300, 0x301}}, []byte("....."), 0, []byte{'a', 0xCC, 0x81, 0xCC, 0x80}, 5, true},
		{&toUnicodeAccents{letter: 'a', accents: []rune{0x300, 0x301, 0x302}}, []byte("....."), 0, []byte{'a', 0xCC, 0x82, 0xCC, 0x81}, 0, false},
	}

	for i, d := range data {
		ok := d.t.write(d.dst, &d.ndst)
		if ok != d.expok {
			t.Errorf("test %d: expected ok=%v, got ok=%v", i, d.expok, ok)
		}
		if !bytes.Equal(d.dst, d.exp) {
			t.Errorf("test %d: expected dst=%s, got dst=%s", i, d.exp, d.dst)
		}
		if d.ndst != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, d.ndst)
		}
	}
}

func TestGetSpecial(t *testing.T) {
	data := []struct {
		src     string
		expls   latexSpecial
		expn    int
		expMore bool
	}{
		{"", latexSpecial{spType: latexSpecialNone}, 0, true},
		{"K", latexSpecial{spType: latexSpecialNone}, 0, false},
		{"oups", latexSpecial{spType: latexSpecialNone}, 0, false},
		{"L", latexSpecial{spType: latexSpecialLetter, utf8: 'Ł'}, 1, true},
		{"L ", latexSpecial{spType: latexSpecialLetter, utf8: 'Ł'}, 2, false},
		{"L a", latexSpecial{spType: latexSpecialLetter, utf8: 'Ł'}, 2, false},
		{"L \\", latexSpecial{spType: latexSpecialLetter, utf8: 'Ł'}, 2, false},
		{"L{}", latexSpecial{spType: latexSpecialLetter, utf8: 'Ł'}, 3, false},
		{"o", latexSpecial{spType: latexSpecialLetter, utf8: 'ø'}, 1, true},
		{"o ", latexSpecial{spType: latexSpecialLetter, utf8: 'ø'}, 2, false},
		{"oe", latexSpecial{spType: latexSpecialLetter, utf8: 'œ'}, 2, true},
		{"oe ", latexSpecial{spType: latexSpecialLetter, utf8: 'œ'}, 3, false},
		{"oe{}", latexSpecial{spType: latexSpecialLetter, utf8: 'œ'}, 4, false},
		{"AE", latexSpecial{spType: latexSpecialLetter, utf8: 'Æ'}, 2, true},
		{"AE ", latexSpecial{spType: latexSpecialLetter, utf8: 'Æ'}, 3, false},
		{"AE{}", latexSpecial{spType: latexSpecialLetter, utf8: 'Æ'}, 4, false},
		{"`", latexSpecial{spType: latexSpecialNonLetterAccent, utf8: 0x300}, 1, false},
		{"` ", latexSpecial{spType: latexSpecialNonLetterAccent, utf8: 0x300}, 1, false},
		{"`{}", latexSpecial{spType: latexSpecialNonLetterAccent, utf8: 0x300}, 1, false},
		{"`{", latexSpecial{spType: latexSpecialNonLetterAccent, utf8: 0x300}, 1, false},
		{"` {", latexSpecial{spType: latexSpecialNonLetterAccent, utf8: 0x300}, 1, false},
		{"c", latexSpecial{spType: latexSpecialLetterAccent, utf8: 0x327}, 1, true},
		{"c ", latexSpecial{spType: latexSpecialLetterAccent, utf8: 0x327}, 2, false},
		{"c{}", latexSpecial{spType: latexSpecialLetterAccent, utf8: 0x327}, 1, false},
	}

	for i, d := range data {
		s, n, more := getSpecial([]byte(d.src))
		if s.spType != d.expls.spType {
			t.Errorf("test %d: expected spType=%v, got spType=%v", i, d.expls.spType, s.spType)
		}
		if s.utf8 != d.expls.utf8 {
			t.Errorf("test %d: expected utf8=%v, got utf8=%v", i, d.expls.utf8, s.utf8)
		}
		if n != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, n)
		}
		if more != d.expMore {
			t.Errorf("test %d: expected more=%v, got more=%v", i, d.expMore, more)
		}
	}
}

func TestGetLetter(t *testing.T) {
	data := []struct {
		src     string
		expl    rune
		expn    int
		expMore bool
	}{
		{"", 0, 0, true},
		{"K", 'K', 1, false},
		{"Kr", 'K', 1, false},
		{"{", 0, 0, true},
		{"{K", 0, 0, true},
		{"{K}", 'K', 3, false},
		{"{K}r", 'K', 3, false},
		{"{Kr", 0, 0, false},
		{"{K ", 0, 0, false},
		{"{K }", 0, 0, false},
		{"{\\L}", 'Ł', 4, false},
		{"{\\L }", 'Ł', 5, false},
		{"{\\L{}}", 'Ł', 6, false},
		{"{\\L", 0, 0, true},
		{"\\L", 0, 0, false},
	}

	for i, d := range data {
		s, n, more := getLetter([]byte(d.src))
		if s != d.expl {
			t.Errorf("test %d: expected letter=%v, got letter=%v", i, d.expl, s)
		}
		if n != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, n)
		}
		if more != d.expMore {
			t.Errorf("test %d: expected more=%v, got more=%v", i, d.expMore, more)
		}
	}
}

func TestToUnicodeAccents_Transform(t *testing.T) {
	data := []struct {
		lendst int    // the length of the destination slice
		src    string // source string
		atEOF  bool   // at EOF
		exp    []byte // expected destination slice
	}{
		{100, "", true, []byte("")},
		{100, "a", true, []byte("a")},
		{100, "\\`a", true, []byte{'a', 0xCC, 0x80}},
		{100, "\\`{a}", true, []byte{'a', 0xCC, 0x80}},
		{100, "\\`{}a", true, []byte{0xCC, 0x80, 'a'}},
		{100, "\\'\\`a", true, []byte{'a', 0xCC, 0x80, 0xCC, 0x81}},
		{100, "\\'\\`", true, []byte{0xCC, 0x80, 0xCC, 0x81}},
		{100, "\\'\\`{", true, []byte{0xCC, 0x80, 0xCC, 0x81, '{'}},
		{100, "\\'\\`{}", true, []byte{0xCC, 0x80, 0xCC, 0x81}},
		{100, "\\'\\`{a", true, []byte{0xCC, 0x80, 0xCC, 0x81, '{', 'a'}},
		{100, "\\'\\`{a}", true, []byte{'a', 0xCC, 0x80, 0xCC, 0x81}},
		{100, "\\'\\c", true, []byte{0xCC, 0xA7, 0xCC, 0x81}},
		{100, "\\'\\c{", true, []byte{0xCC, 0xA7, 0xCC, 0x81, '{'}},
		{100, "\\'\\c{}", true, []byte{0xCC, 0xA7, 0xCC, 0x81}},
		{100, "\\'\\c{c", true, []byte{0xCC, 0xA7, 0xCC, 0x81, '{', 'c'}},
		{100, "\\'\\c{c}", true, []byte{'c', 0xCC, 0xA7, 0xCC, 0x81}},
		{100, "\\`\\L", true, []byte{0xC5, 0x81, 0xCC, 0x80}},
		{100, "{\\`\\L}", true, []byte{0xC5, 0x81, 0xCC, 0x80}},
	}

	for i, d := range data {
		lat := &toUnicodeAccents{}
		dst := make([]byte, d.lendst)
		nDst, nSrc, err := lat.Transform(dst, []byte(d.src), d.atEOF)
		if err != nil {
			t.Errorf("test %d: expected err=nil, got err=%v", i, err)
		}
		dst = dst[:nDst]
		if !bytes.Equal(dst, d.exp) {
			t.Errorf("test %d: expected dst=%v, got dst=%v", i, d.exp, dst)
		}
		if nSrc != len(d.src) {
			t.Errorf("test %d: expected nSrc=%d, got nSrc=%d", i, len(d.src), nSrc)
		}
	}
}
