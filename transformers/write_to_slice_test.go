package transformers

import (
	"bytes"
	"testing"
)

func TestWriteByte(t *testing.T) {
	data := []struct {
		dst    []byte // destination slice
		src    byte   // source byte to write
		expdst []byte // expected destination slice after writing
		ndst   int    // nuber of bytes already written to the destination slice
		expn   int    // expected number of bytes written to the destination slice after writing
		expok  bool   // expected return value of the write operation
	}{
		{[]byte{0, 0, 0, 0}, 0x01, []byte{0x01, 0, 0, 0}, 0, 1, true},
		{[]byte{0, 0, 0, 0}, 0x01, []byte{0, 0x01, 0, 0}, 1, 2, true},
		{[]byte{0, 0, 0, 0}, 0x01, []byte{0, 0, 0x01, 0}, 2, 3, true},
		{[]byte{0, 0, 0, 0}, 0x01, []byte{0, 0, 0, 0x01}, 3, 4, true},
		{[]byte{0, 0, 0, 0}, 0x01, []byte{0, 0, 0, 0}, 4, 4, false},
	}

	for i, d := range data {
		ok := writeByte(d.dst, d.src, &d.ndst)
		if ok != d.expok {
			t.Errorf("test %d: expected ok=%v, got ok=%v", i, d.expok, ok)
		}
		if !bytes.Equal(d.dst, d.expdst) {
			t.Errorf("test %d: expected dst=%v, got dst=%v", i, d.expdst, d.dst)
		}
		if d.ndst != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, d.ndst)
		}
	}
}

func TestWriteBytes(t *testing.T) {
	data := []struct {
		dst    []byte // destination slice
		src    []byte // source byte to write
		expdst []byte // expected destination slice after writing
		ndst   int    // nuber of bytes already written to the destination slice
		expn   int    // expected number of bytes written to the destination slice after writing
		expok  bool   // expected return value of the write operation
	}{
		{[]byte{0, 0, 0, 0}, []byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03, 0}, 0, 3, true},
		{[]byte{0, 0, 0, 0}, []byte{0x01, 0x02, 0x03}, []byte{0, 0x01, 0x02, 0x03}, 1, 4, true},
		{[]byte{0, 0, 0, 0}, []byte{0x01, 0x02, 0x03}, []byte{0, 0, 0x01, 0x02}, 2, 2, false},
		{[]byte{0, 0, 0, 0}, []byte{0x01, 0x02, 0x03}, []byte{0, 0, 0, 0x01}, 3, 3, false},
		{[]byte{0, 0, 0, 0}, []byte{0x01, 0x02, 0x03}, []byte{0, 0, 0, 0}, 4, 4, false},
	}

	for i, d := range data {
		ok := writeBytes(d.dst, d.src, &d.ndst)
		if ok != d.expok {
			t.Errorf("test %d: expected ok=%v, got ok=%v", i, d.expok, ok)
		}
		if !bytes.Equal(d.dst, d.expdst) {
			t.Errorf("test %d: expected dst=%v, got dst=%v", i, d.expdst, d.dst)
		}
		if d.ndst != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, d.ndst)
		}
	}
}

func TestWriteString(t *testing.T) {
	data := []struct {
		dst    []byte // destination slice
		src    string // source byte to write
		expdst []byte // expected destination slice after writing
		ndst   int    // nuber of bytes already written to the destination slice
		expn   int    // expected number of bytes written to the destination slice after writing
		expok  bool   // expected return value of the write operation
	}{
		{[]byte{0, 0, 0, 0}, "abc", []byte{'a', 'b', 'c', 0}, 0, 3, true},
		{[]byte{0, 0, 0, 0}, "abc", []byte{0, 'a', 'b', 'c'}, 1, 4, true},
		{[]byte{0, 0, 0, 0}, "abc", []byte{0, 0, 'a', 'b'}, 2, 2, false},
		{[]byte{0, 0, 0, 0}, "abc", []byte{0, 0, 0, 'a'}, 3, 3, false},
		{[]byte{0, 0, 0, 0}, "abc", []byte{0, 0, 0, 0}, 4, 4, false},
	}

	for i, d := range data {
		ok := writeString(d.dst, d.src, &d.ndst)
		if ok != d.expok {
			t.Errorf("test %d: expected ok=%v, got ok=%v", i, d.expok, ok)
		}
		if !bytes.Equal(d.dst, d.expdst) {
			t.Errorf("test %d: expected dst=%v, got dst=%v", i, d.expdst, d.dst)
		}
		if d.ndst != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, d.ndst)
		}
	}
}

func TestWriteRune(t *testing.T) {
	data := []struct {
		dst    []byte // destination slice
		src    rune   // source byte to write
		expdst []byte // expected destination slice after writing
		ndst   int    // nuber of bytes already written to the destination slice
		expn   int    // expected number of bytes written to the destination slice after writing
		expok  bool   // expected return value of the write operation
	}{
		{[]byte{0, 0, 0, 0}, 'a', []byte{'a', 0, 0, 0}, 0, 1, true},
		{[]byte{0, 0, 0, 0}, 'a', []byte{0, 'a', 0, 0}, 1, 2, true},
		{[]byte{0, 0, 0, 0}, 'a', []byte{0, 0, 'a', 0}, 2, 3, true},
		{[]byte{0, 0, 0, 0}, 'a', []byte{0, 0, 0, 'a'}, 3, 4, true},
		{[]byte{0, 0, 0, 0}, 'a', []byte{0, 0, 0, 0}, 4, 4, false},
		{[]byte{0, 0, 0, 0}, 'α', []byte{0xCE, 0xB1, 0, 0}, 0, 2, true},
		// U+1F642 -> 0xF0, 0x9F, 0x99, 0x82
		{[]byte{0, 0, 0, 0}, '🙂', []byte{0xF0, 0x9F, 0x99, 0x82}, 0, 4, true},
		{[]byte{0, 0, 0, 0}, '🙂', []byte{0, 0xF0, 0x9F, 0x99}, 1, 1, false},
	}

	for i, d := range data {
		ok := writeLaTeXRune(d.dst, d.src, &d.ndst)
		if ok != d.expok {
			t.Errorf("test %d: expected ok=%v, got ok=%v", i, d.expok, ok)
		}
		if !bytes.Equal(d.dst, d.expdst) {
			t.Errorf("test %d: expected dst=%v, got dst=%v", i, d.expdst, d.dst)
		}
		if d.ndst != d.expn {
			t.Errorf("test %d: expected n=%d, got n=%d", i, d.expn, d.ndst)
		}
	}
}

func TestIsLatinByte(t *testing.T) {
	data := []struct {
		r   byte // rune to test
		exp bool // expected return value
	}{
		{'a', true},
		{'A', true},
		{'z', true},
		{'Z', true},
		{'0', false},
		{'9', false},
		{' ', false},
		{'\\', false},
	}

	for i, d := range data {
		if isLatinByte(d.r) != d.exp {
			t.Errorf("test %d: expected %v, got %v", i, d.exp, !d.exp)
		}
	}
}
