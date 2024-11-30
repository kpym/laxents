package transformers

import (
	"unicode/utf8"
)

// bytestr is a type constraint for []byte and string, used for functions
// that operate generically on these types.
type bytestr interface {
	~[]byte | ~string
}

// writeBytes writes src to dst and increments nDst by the number of bytes written
// it returns true if the whole slice was written
func write[T bytestr](dst []byte, src T, nDst *int) bool {
	n := copy(dst[*nDst:], src)
	if n < len(src) {
		return false
	}
	*nDst += n
	return true
}

// writeByte writes b to dst and increments nDst by 1
func writeByte(dst []byte, b byte, nDst *int) bool {
	if *nDst < len(dst) {
		dst[*nDst] = b
		*nDst++
		return true
	}
	return false
}

// buf is a buffer to store the UTF-8 encoding of a rune
var buf = make([]byte, utf8.UTFMax)

// writeLaTeXRune writes r to dst and increments nDst by the number of bytes written
// it returns true if the whole rune was written
func writeRune(dst []byte, r rune, nDst *int) bool {
	n := utf8.EncodeRune(buf, r)
	return write(dst, buf[:n], nDst)
}

type byterune interface {
	~byte | ~rune
}

// isLatinByte returns true if r is a latin letter
func isLatin[T byterune](r T) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}
