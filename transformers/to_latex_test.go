package transformers

import (
	"bytes"
	"testing"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func TestToLaTeXAccentsReset(t *testing.T) {
	lat := &toLaTeXAccents{letter: 'a', accents: []rune{'`'}}
	lat.Reset()
	if lat.letter != 0 {
		t.Errorf("expected letter=0, got letter=%v", lat.letter)
	}
	if len(lat.accents) != 0 {
		t.Errorf("expected accents=nil, got accents=%v", lat.accents)
	}
}

func TestWriteLaTeXRune(t *testing.T) {
	data := []struct {
		r      rune   // rune
		dst    []byte // destination slice
		ndst   int    // number of bytes already written to the destination slice
		expdst []byte // expected destination slice after writing
		expres bool   // expected return value of the write operation
	}{
		{'a', []byte("....."), 0, []byte("a...."), true},
		{'a', []byte("....."), 4, []byte("....a"), true},
		{'a', []byte("....."), 5, []byte("....."), false},
		{'Ł', []byte("....."), 0, []byte("{\\L}."), true},
		{'Ł', []byte("....."), 4, []byte("....{"), false},
	}

	for i, d := range data {
		ok := writeLaTeXRune(d.dst, d.r, &d.ndst)
		if ok != d.expres {
			t.Errorf("test %d: expected ok=%v, got ok=%v", i, d.expres, ok)
		}
		if !bytes.Equal(d.dst, d.expdst) {
			t.Errorf("test %d: expected dst=%s, got dst=%s", i, d.expdst, d.dst)
		}
	}
}

func TestWriteLaTeXLetter(t *testing.T) {
	data := []struct {
		letter  rune   // letter rune
		dst     []byte // destination slice
		ndst    int    // number of bytes already written to the destination slice
		inGroup bool   // letter is in a group
		exp     []byte // expected destination slice after writing
		expn    int    // expected number of bytes written to the destination slice after writing
		expok   bool   // expected return value of the write operation
	}{
		{'a', []byte("....."), 0, false, []byte("a...."), 1, true},
		{'a', []byte("....."), 4, false, []byte("....a"), 5, true},
		{'a', []byte("....."), 5, false, []byte("....."), 5, false},
		{'a', []byte("....."), 0, true, []byte("{a}.."), 3, true},
		{'a', []byte("....."), 2, true, []byte("..{a}"), 5, true},
		{'a', []byte("....."), 3, true, []byte("...{a"), 3, false},
		{'a', []byte("....."), 4, true, []byte("....{"), 4, false},
		{'a', []byte("....."), 5, true, []byte("....."), 5, false},
	}

	for i, d := range data {
		lat := &toLaTeXAccents{letter: d.letter}
		ok := lat.writeLaTeXLetter(d.dst, &d.ndst, d.inGroup)
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

func TestWriteLaTeXAccent(t *testing.T) {
	data := []struct {
		dst    []byte // destination slice
		acc    []rune // accent rune
		letter rune   // letter rune
		exp    []byte // expected destination slice after writing
		ndst   int    // number of bytes already written to the destination slice
		expn   int    // expected number of bytes written to the destination slice after writing
		expok  bool   // expected return value of the write operation
	}{
		{[]byte("....."), []rune("`"), 'a', []byte("\\`a.."), 0, 3, true},
		{[]byte("....."), []rune("`"), 'a', []byte(".\\`a."), 1, 4, true},
		{[]byte("....."), []rune("`"), 'a', []byte("..\\`a"), 2, 5, true},
		{[]byte("....."), []rune("`"), 'a', []byte("...\\`"), 3, 3, false},
		{[]byte("....."), []rune("`"), 'a', []byte("....\\"), 4, 4, false},
		{[]byte("....."), []rune("`"), 'a', []byte("....."), 5, 5, false},
		{[]byte("....."), []rune("c"), 'c', []byte("\\c{c}"), 0, 5, true},
		{[]byte("....."), []rune("c"), 'c', []byte(".\\c{c"), 1, 1, false},
		{[]byte("....."), []rune("c"), 'c', []byte("..\\c{"), 2, 2, false},
		{[]byte("....."), []rune("c"), 'c', []byte("...\\c"), 3, 3, false},
		{[]byte("....."), []rune("c"), 'c', []byte("....\\"), 4, 4, false},
		{[]byte("....."), []rune("c"), 'c', []byte("....."), 5, 5, false},
		{[]byte("......."), []rune("`'"), 'a', []byte("\\'\\`a.."), 0, 5, true},
		{[]byte("......."), []rune("`'^"), 'a', []byte("\\^\\'\\`a"), 0, 7, true},
	}

	for i, d := range data {
		lat := &toLaTeXAccents{accents: d.acc, letter: d.letter}
		ok := lat.writeLaTeXAccent(d.dst, &d.ndst)
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

func TestToLaTeXAccents(t *testing.T) {
	data := []struct {
		src string // source string
		exp string // expected destination string after writing
	}{
		{"a", "a"},
		{"ç", "\\c{c}"},
		{"é", "\\'e"},
		{"bêtâ", "b\\^et\\^a"},
		{"ḵ", "\\b{k}"},
		{"Ceci est œuf", "Ceci est {\\oe}uf"},
	}

	var buf bytes.Buffer
	w := transform.NewWriter(&buf, ToLaTeXAccents())

	for i, d := range data {
		buf.Reset()
		// convert d.src to NFD (required by ToLaTeXAccents)
		src := norm.NFD.String(d.src)
		_, err := w.Write([]byte(src))
		if err != nil {
			t.Errorf("test %d: unexpected error: %v", i, err)
		}
		err = w.Close()
		if err != nil {
			t.Errorf("test %d: unexpected error: %v", i, err)
		}
		if buf.String() != d.exp {
			t.Errorf("test %d: expected dst=%v, got dst=%v", i, d.exp, buf.String())
		}
	}
}
