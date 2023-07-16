// Package clockwork implements Clockwork Base32 encoding as specified by https://gist.github.com/szktty/228f85794e4187882a77734c89c384a8
package clockwork

import (
	"io"
	"strconv"
)

// Base32 is Clockwork Base32 encoding.
var Base32 = NewEncoding()

/*
 * Encodings
 */

// An Encoding is a radix 32 encoding/decoding scheme.
type Encoding struct {
	encode    [32]byte
	decodeMap [256]byte
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
		decodeMap: [256]byte{
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 0-9 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 10-19 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 20-29 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 30-39 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0, 1, /* 40-49 */
			2, 3, 4, 5, 6, 7, 8, 9, 0, 0xFF, /* 50-59 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 10, 11, 12, 13, 14, /* 60-69 */
			15, 16, 17, 1, 18, 19, 1, 20, 21, 0, /* 70-79 */
			22, 23, 24, 25, 26, 0xFF, 27, 28, 29, 30, /* 80-89 */
			31, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 10, 11, 12, /* 90-99 */
			13, 14, 15, 16, 17, 1, 18, 19, 1, 20, /* 100-109 */
			21, 0, 22, 23, 24, 25, 26, 0xFF, 27, 28, /* 110-119 */
			29, 30, 31, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 120-129 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 130-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 140-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 150-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 160-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 170-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 180-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 190-109 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 200-209 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 210-209 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 220-209 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 230-209 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 240-209 */
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, /* 250-256 */
		},
	}
}

// Encode encodes src using the encoding enc, writing
// EncodedLen(len(src)) bytes to dst.
func (enc *Encoding) Encode(dst, src []byte) {
	for len(src) >= 5 {
		// Unpack 8x 5-bit source blocks into a 5 byte
		// destination quantum
		val := uint64(src[0])<<32 | uint64(src[1])<<24 | uint64(src[2])<<16 | uint64(src[3])<<8 | uint64(src[4])
		dst[0] = enc.encode[(val>>35)&0x1F]
		dst[1] = enc.encode[(val>>30)&0x1F]
		dst[2] = enc.encode[(val>>25)&0x1F]
		dst[3] = enc.encode[(val>>20)&0x1F]
		dst[4] = enc.encode[(val>>15)&0x1F]
		dst[5] = enc.encode[(val>>10)&0x1F]
		dst[6] = enc.encode[(val>>5)&0x1F]
		dst[7] = enc.encode[(val>>0)&0x1F]
		src = src[5:]
		dst = dst[8:]
	}

	// Add the remaining small block
	if len(src) > 0 {
		var val uint64
		switch len(src) {
		default:
			val |= uint64(src[4])
			fallthrough
		case 4:
			val |= uint64(src[3]) << 8
			fallthrough
		case 3:
			val |= uint64(src[2]) << 16
			fallthrough
		case 2:
			val |= uint64(src[1]) << 24
			fallthrough
		case 1:
			val |= uint64(src[0]) << 32
		}

		// Encode 5-bit blocks using the base32 alphabet
		size := uint(len(dst))
		if size >= 8 {
			size = 8
		}
		for i := uint(0); i < size; i++ {
			dst[i] = enc.encode[(val>>(35-5*i))&31]
		}
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
	copy(e.buf[:], p)
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

/*
 * Decoder
 */

// CorruptInputError is a decoding error.
type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal clockwork base32 data at input byte " + strconv.FormatInt(int64(e), 10)
}

// decode is like Decode but returns an additional 'end' value, which
// indicates if end-of-message padding was encountered and thus any
// additional data is an error. This method assumes that src has been
// stripped of all supported whitespace ('\r' and '\n').
func (enc *Encoding) decode(dst, src []byte) (n int, err error) {
	// Lift the nil check outside of the loop.
	_ = enc.decodeMap

	olen := len(src)

	// Decode in 8-byte chunks
	for len(src) >= 8 {
		var dbuf [8]byte
		for j := 0; j < len(dbuf); j++ {
			dbuf[j] = enc.decodeMap[src[j]]
			if dbuf[j] == 0xFF {
				return n, CorruptInputError(olen - len(src) + j)
			}
		}
		src = src[8:]

		// Pack 8x 5-bit source blocks into 5 byte destination
		// quantum
		val := uint64(dbuf[0])<<35 |
			uint64(dbuf[1])<<30 |
			uint64(dbuf[2])<<25 |
			uint64(dbuf[3])<<20 |
			uint64(dbuf[4])<<15 |
			uint64(dbuf[5])<<10 |
			uint64(dbuf[6])<<5 |
			uint64(dbuf[7])
		dst[0] = byte(val >> 32)
		dst[1] = byte(val >> 24)
		dst[2] = byte(val >> 16)
		dst[3] = byte(val >> 8)
		dst[4] = byte(val)
		n += 5
		dst = dst[5:]
	}

	// Add the remaining small block
	if len(src) > 0 {
		// Decode quantum using the base32 alphabet
		var dbuf [8]byte
		for j := 0; j < len(src); j++ {
			in := src[j]
			dbuf[j] = enc.decodeMap[in]
			if dbuf[j] == 0xFF {
				return n, CorruptInputError(olen - len(src) - 1)
			}
		}

		// Pack 8x 5-bit source blocks into 5 byte destination
		// quantum
		val := uint64(dbuf[0])<<35 |
			uint64(dbuf[1])<<30 |
			uint64(dbuf[2])<<25 |
			uint64(dbuf[3])<<20 |
			uint64(dbuf[4])<<15 |
			uint64(dbuf[5])<<10 |
			uint64(dbuf[6])<<5 |
			uint64(dbuf[7])
		switch len(src) {
		case 8:
			dst[4] = byte(val)
			n++
			fallthrough
		case 7:
			dst[3] = byte(val >> 8)
			n++
			fallthrough
		case 6, 5:
			// dbuf[5] might be padding
			dst[2] = byte(val >> 16)
			n++
			fallthrough
		case 4:
			dst[1] = byte(val >> 24)
			n++
			fallthrough
		case 3, 2:
			// dbuf[2] might be padding
			dst[0] = byte(val >> 32)
			n++
		}
	}
	return n, nil
}

// Decode decodes src using the encoding enc. It writes at most
// DecodedLen(len(src)) bytes to dst and returns the number of bytes
// written. If src contains invalid base32 data, it will return the
// number of bytes successfully written and CorruptInputError.
func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
	return enc.decode(dst, src)
}

// DecodeString returns the bytes represented by the base32 string s.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	buf := []byte(s)
	n, err := enc.decode(buf, buf)
	return buf[:n], err
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base32-encoded data.
func (enc *Encoding) DecodedLen(n int) int {
	return n * 5 / 8
}

type decoder struct {
	err    error
	enc    *Encoding
	r      io.Reader
	buf    [1024]byte // leftover input
	nbuf   int
	out    []byte // leftover decoded output
	outbuf [1024 / 8 * 5]byte
}

func readEncodedData(r io.Reader, buf []byte) (n int, err error) {
	for n < 1 && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	return
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// based on https://github.com/golang/go/blob/ba9e10889976025ee1d027db6b1cad383ec56de8/src/encoding/base32/base32.go#L410

	// Use leftover decoded output from last read.
	if len(d.out) > 0 {
		n = copy(p, d.out)
		d.out = d.out[n:]
		if len(d.out) == 0 {
			return n, d.err
		}
		return n, nil
	}

	if d.err != nil {
		return 0, d.err
	}

	// Read a chunk.
	nn := len(p) / 5 * 8
	if nn < 8 {
		nn = 8
	}
	if nn > len(d.buf) {
		nn = len(d.buf)
	}
	nn, d.err = readEncodedData(d.r, d.buf[d.nbuf:nn])
	d.nbuf += nn

	// Decode chunk into p, or d.out and then p if p is too small.
	nr := d.nbuf
	nw := d.enc.DecodedLen(d.nbuf)
	if nw > len(p) {
		nw, err = d.enc.decode(d.outbuf[0:], d.buf[0:nr])
		d.out = d.outbuf[0:nw]
		n = copy(p, d.out)
		d.out = d.out[n:]
	} else {
		n, err = d.enc.decode(p, d.buf[0:nr])
	}
	d.nbuf -= nr

	for i := 0; i < d.nbuf; i++ {
		d.buf[i] = d.buf[i+nr]
	}

	if err != nil && (d.err == nil || d.err == io.EOF) {
		d.err = err
	}

	if len(d.out) > 0 {
		// We cannot return all the decoded bytes to the caller in this
		// invocation of Read, so we return a nil error to ensure that Read
		// will be called again.  The error stored in d.err, if any, will be
		// returned with the last set of decoded bytes.
		return n, nil
	}

	return n, d.err
}

// NewDecoder constructs a new base32 stream decoder.
func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	return &decoder{enc: enc, r: r}
}
