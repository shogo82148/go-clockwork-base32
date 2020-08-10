// Package clockwork implements Clockwork Base32 encoding as specified by https://gist.github.com/szktty/228f85794e4187882a77734c89c384a8
package clockwork

import "io"

/*
 * Encodings
 */

// An Encoding is a radix 32 encoding/decoding scheme.
type Encoding struct {
	encode    [32]byte
	decodeMap [256]int8
}

// NewEncoding returns a new Encoding.
func NewEncoding() *Encoding {
	return &Encoding{
		// https://github.com/szktty/go-clockwork-base32/blob/c2cac4daa7ad2045089b943b377b12ac57e3254e/base32.go#L61-L66
		encode: [32]byte{
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K',
			'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X',
			'Y', 'Z',
		},
		// https://github.com/szktty/go-clockwork-base32/blob/c2cac4daa7ad2045089b943b377b12ac57e3254e/base32.go#L68-L95
		decodeMap: [256]int8{
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 0-9 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 10-19 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 20-29 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 30-39 */
			-1, -1, -1, -1, -1, -1, -1, -1, 0, 1, /* 40-49 */
			2, 3, 4, 5, 6, 7, 8, 9, 0, -1, /* 50-59 */
			-1, -1, -1, -1, -1, 10, 11, 12, 13, 14, /* 60-69 */
			15, 16, 17, 1, 18, 19, 1, 20, 21, 0, /* 70-79 */
			22, 23, 24, 25, 26, -2, 27, 28, 29, 30, /* 80-89 */
			31, -1, -1, -1, -1, -1, -1, 10, 11, 12, /* 90-99 */
			13, 14, 15, 16, 17, 1, 18, 19, 1, 20, /* 100-109 */
			21, 0, 22, 23, 24, 25, 26, -1, 27, 28, /* 110-119 */
			29, 30, 31, -1, -1, -1, -1, -1, -1, -1, /* 120-129 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 130-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 140-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 150-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 160-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 170-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 180-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 190-109 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 200-209 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 210-209 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 220-209 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 230-209 */
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, /* 240-209 */
			-1, -1, -1, -1, -1, -1, /* 250-256 */
		},
	}
}

// Encode encodes src using the encoding enc, writing
// EncodedLen(len(src)) bytes to dst.
func (enc *Encoding) Encode(dst, src []byte) {
	// based on https://github.com/golang/go/blob/ba9e10889976025ee1d027db6b1cad383ec56de8/src/encoding/base32/base32.go#L93

	for len(src) > 0 {
		var b [8]byte

		// Unpack 8x 5-bit source blocks into a 5 byte
		// destination quantum
		switch len(src) {
		default:
			b[7] = src[4] & 0x1F
			b[6] = src[4] >> 5
			fallthrough
		case 4:
			b[6] |= (src[3] << 3) & 0x1F
			b[5] = (src[3] >> 2) & 0x1F
			b[4] = src[3] >> 7
			fallthrough
		case 3:
			b[4] |= (src[2] << 1) & 0x1F
			b[3] = (src[2] >> 4) & 0x1F
			fallthrough
		case 2:
			b[3] |= (src[1] << 4) & 0x1F
			b[2] = (src[1] >> 1) & 0x1F
			b[1] = (src[1] >> 6) & 0x1F
			fallthrough
		case 1:
			b[1] |= (src[0] << 2) & 0x1F
			b[0] = src[0] >> 3
		}

		// Encode 5-bit blocks using the base32 alphabet
		size := len(dst)
		if size >= 8 {
			// Common case, unrolled for extra performance
			dst[0] = enc.encode[b[0]&31]
			dst[1] = enc.encode[b[1]&31]
			dst[2] = enc.encode[b[2]&31]
			dst[3] = enc.encode[b[3]&31]
			dst[4] = enc.encode[b[4]&31]
			dst[5] = enc.encode[b[5]&31]
			dst[6] = enc.encode[b[6]&31]
			dst[7] = enc.encode[b[7]&31]
		} else {
			for i := 0; i < size; i++ {
				dst[i] = enc.encode[b[i]&31]
			}
		}

		if len(src) < 5 {
			break
		}
		src = src[5:]
		dst = dst[8:]
	}
}

// EncodeToString returns the base32 encoding of src.
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}

// EncodedLen returns the length in bytes of the base32 encoding
// of an input buffer of length n.
func (enc *Encoding) EncodedLen(n int) int {
	return (n*8 + 4) / 5
}

type encoder struct {
	err  error
	enc  *Encoding
	w    io.Writer
	buf  [5]byte    // buffered data waiting to be encoded
	nbuf int        // number of bytes in buf
	out  [1024]byte // output buffer
}

func (e *encoder) Write(p []byte) (n int, err error) {
	// based on https://github.com/golang/go/blob/ba9e10889976025ee1d027db6b1cad383ec56de8/src/encoding/base32/base32.go#L184

	if e.err != nil {
		return 0, e.err
	}

	// Leading fringe.
	if e.nbuf > 0 {
		var i int
		for i = 0; i < len(p) && e.nbuf < 5; i++ {
			e.buf[e.nbuf] = p[i]
			e.nbuf++
		}
		n += i
		p = p[i:]
		if e.nbuf < 5 {
			return
		}
		e.enc.Encode(e.out[0:], e.buf[0:])
		if _, e.err = e.w.Write(e.out[0:8]); e.err != nil {
			return n, e.err
		}
		e.nbuf = 0
	}

	// Large interior chunks.
	for len(p) >= 5 {
		nn := len(e.out) / 8 * 5
		if nn > len(p) {
			nn = len(p)
			nn -= nn % 5
		}
		e.enc.Encode(e.out[0:], p[0:nn])
		if _, e.err = e.w.Write(e.out[0 : nn/5*8]); e.err != nil {
			return n, e.err
		}
		n += nn
		p = p[nn:]
	}

	// Trailing fringe.
	for i := 0; i < len(p); i++ {
		e.buf[i] = p[i]
	}
	e.nbuf = len(p)
	n += len(p)
	return
}

// Close flushes any pending output from the encoder.
// It is an error to call Write after calling Close.
func (e *encoder) Close() error {
	// If there's anything left in the buffer, flush it out
	if e.err == nil && e.nbuf > 0 {
		e.enc.Encode(e.out[0:], e.buf[0:e.nbuf])
		encodedLen := e.enc.EncodedLen(e.nbuf)
		e.nbuf = 0
		_, e.err = e.w.Write(e.out[0:encodedLen])
	}
	return e.err
}

// NewEncoder returns a new base32 stream encoder. Data written to
// the returned writer will be encoded using enc and then written to w.
// Base32 encodings operate in 5-byte blocks; when finished
// writing, the caller must Close the returned encoder to flush any
// partially written blocks.
func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	return &encoder{enc: enc, w: w}
}
