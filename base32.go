// Package clockwork implements Clockwork Base32 encoding as specified by https://gist.github.com/szktty/228f85794e4187882a77734c89c384a8
package clockwork

import (
	"io"
	"strconv"
)

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
	// based on https://github.com/golang/go/blob/ba9e10889976025ee1d027db6b1cad383ec56de8/src/encoding/base32/base32.go#L277

	// Lift the nil check outside of the loop.
	_ = enc.decodeMap

	dsti := 0
	olen := len(src)

	for len(src) > 0 {
		// Decode quantum using the base32 alphabet
		var dbuf [8]byte
		dlen := len(src)
		if dlen > 8 {
			dlen = 8
		}

		for j := 0; j < dlen; j++ {
			in := src[0]
			src = src[1:]
			dbuf[j] = enc.decodeMap[in]
			if dbuf[j] == 0xFF {
				return n, CorruptInputError(olen - len(src) - 1)
			}
		}

		// Pack 8x 5-bit source blocks into 5 byte destination
		// quantum
		switch dlen {
		case 8:
			dst[dsti+4] = dbuf[6]<<5 | dbuf[7]
			n++
			fallthrough
		case 7:
			dst[dsti+3] = dbuf[4]<<7 | dbuf[5]<<2 | dbuf[6]>>3
			n++
			fallthrough
		case 6, 5:
			// dbuf[5] might be padding
			dst[dsti+2] = dbuf[3]<<4 | dbuf[4]>>1
			n++
			fallthrough
		case 4:
			dst[dsti+1] = dbuf[1]<<6 | dbuf[2]<<1 | dbuf[3]>>4
			n++
			fallthrough
		case 3, 2:
			// dbuf[2] might be padding
			dst[dsti+0] = dbuf[0]<<3 | dbuf[1]>>2
			n++
		}
		dsti += 5
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
	end    bool       // saw end of message
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
