package transformers

import (
	"unicode/utf8"
)

// writeByte writes b to dst and increments nDst by 1
func writeByte(dst []byte, b byte, nDst *int) bool {
	if *nDst < len(dst) {
		dst[*nDst] = b
		*nDst++
		return true
	}
	return false
}

// writeBytes writes src to dst and increments nDst by the number of bytes written
// it returns true if the whole slice was written
func writeBytes(dst, src []byte, nDst *int) bool {
	n := copy(dst[*nDst:], src)
	if n < len(src) {
		return false
	}
	*nDst += n
	return true
}

// writeString writes s to dst and increments nDst by the number of bytes written
// it returns true if the whole string was written
func writeString(dst []byte, src string, nDst *int) bool {
	n := copy(dst[*nDst:], src)
	if n < len(src) {
		return false
	}
	*nDst += n
	return true
}

// buf is a buffer to store the UTF-8 encoding of a rune
var buf = make([]byte, utf8.UTFMax)

// writeLaTeXRune writes r to dst and increments nDst by the number of bytes written
// it returns true if the whole rune was written
func writeRune(dst []byte, r rune, nDst *int) bool {
	n := utf8.EncodeRune(buf, r)
	return writeBytes(dst, buf[:n], nDst)
}

// isLatinByte returns true if r is a latin letter
func isLatinByte(r byte) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

// isLatinByte returns true if r is a latin letter
func isLatinRune(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}
