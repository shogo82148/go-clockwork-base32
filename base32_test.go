package clockwork

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type testCase struct {
	plain   string
	encoded string
}

var testCasesEncode = []testCase{
	// from https://github.com/szktty/go-clockwork-base32/blob/c2cac4daa7ad2045089b943b377b12ac57e3254e/base32_test.go#L36-L44
	{"foobar", "CSQPYRK1E8"},
	{"Hello, world!", "91JPRV3F5GG7EVVJDHJ22"},
	{
		"The quick brown fox jumps over the lazy dog.",
		"AHM6A83HENMP6TS0C9S6YXVE41K6YY10D9TPTW3K41QQCSBJ41T6GS90DHGQMY90CHQPEBG",
	},
	{
		"Wow, it really works!",
		"AXQQEB10D5T20WK5C5P6RY90EXQQ4TVK44",
	},

	// from https://github.com/shiguredo/erlang-base32/blob/0cc88a702ce1d8ca345e516a05a9a85f7f23a718/test/base32_clockwork_test.erl#L7-L18
	{"", ""},
	{"f", "CR"},
	{"fo", "CSQG"},
	{"foo", "CSQPY"},
	{"foob", "CSQPYRG"},
	{"fooba", "CSQPYRK1"},
	{"foobar", "CSQPYRK1E8"},
	{
		"\x01\xdd\x3e\x62\xfe\x15\x4e\xd7\x2b\x6d\x2d\x24\x39\x74\x66\x9d",
		"07EKWRQY2N7DEAVD5MJ3JX36KM",
	},
	{
		"Wow, it really works!",
		"AXQQEB10D5T20WK5C5P6RY90EXQQ4TVK44",
	},
}

var bigtest = testCase{
	"The quick brown fox jumps over the lazy dog.",
	"AHM6A83HENMP6TS0C9S6YXVE41K6YY10D9TPTW3K41QQCSBJ41T6GS90DHGQMY90CHQPEBG",
}

func TestEncode(t *testing.T) {
	enc := NewEncoding()
	for _, testCase := range testCasesEncode {
		plain := []byte(testCase.plain)
		encoded := make([]byte, enc.EncodedLen(len(plain)))
		enc.Encode(encoded, plain)
		if !bytes.Equal(encoded, []byte(testCase.encoded)) {
			t.Errorf("encoded %q, expected %q, actual %q\n",
				testCase.plain, testCase.encoded, encoded)
		}
	}
}

func TestEncoder(t *testing.T) {
	enc := NewEncoding()
	for _, testCase := range testCasesEncode {
		var buf bytes.Buffer
		w := NewEncoder(enc, &buf)
		if _, err := w.Write([]byte(testCase.plain)); err != nil {
			t.Errorf("error while encoding %q: %v", testCase.plain, err)
			continue
		}
		if err := w.Close(); err != nil {
			t.Errorf("error while encoding %q: %v", testCase.plain, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), []byte(testCase.encoded)) {
			t.Errorf("encoded %q, expected %q, actual %q\n",
				testCase.plain, testCase.encoded, buf.Bytes())
		}
	}
}

func TestEncoder_Buffering(t *testing.T) {
	input := []byte(bigtest.plain)
	for bs := 1; bs <= 24; bs++ {
		bb := &strings.Builder{}
		encoder := NewEncoder(NewEncoding(), bb)
		for pos := 0; pos < len(input); pos += bs {
			end := pos + bs
			if end > len(input) {
				end = len(input)
			}
			n, err := encoder.Write(input[pos:end])
			if err != nil {
				t.Errorf("error while encoding %q: %v", bigtest.plain, err)
			}
			if n != end-pos {
				t.Errorf("unexpected wrote length: got %d, want %d", n, end-pos)
			}
		}
		if err := encoder.Close(); err != nil {
			t.Errorf("error while encoding %q: %v", bigtest.plain, err)
		}
	}
}

var testCasesDecode = []testCase{
	// from https://github.com/szktty/go-clockwork-base32/blob/c2cac4daa7ad2045089b943b377b12ac57e3254e/base32_test.go#L36-L44
	{"foobar", "CSQPYRK1E8"},
	{"Hello, world!", "91JPRV3F5GG7EVVJDHJ22"},
	{
		"The quick brown fox jumps over the lazy dog.",
		"AHM6A83HENMP6TS0C9S6YXVE41K6YY10D9TPTW3K41QQCSBJ41T6GS90DHGQMY90CHQPEBG",
	},
	{
		"Wow, it really works!",
		"AXQQEB10D5T20WK5C5P6RY90EXQQ4TVK44",
	},

	// from https://github.com/shiguredo/erlang-base32/blob/0cc88a702ce1d8ca345e516a05a9a85f7f23a718/test/base32_clockwork_test.erl#L20-L31
	{"", ""},
	{"f", "CR"},
	{"fo", "CSQG"},
	{"foo", "CSQPY"},
	{"foob", "CSQPYRG"},
	{"fooba", "CSQPYRK1"},
	{"foobar", "CSQPYRK1E8"},
	{
		"\x01\xdd\x3e\x62\xfe\x15\x4e\xd7\x2b\x6d\x2d\x24\x39\x74\x66\x9d",
		"07EKWRQY2N7DEAVD5MJ3JX36KM",
	},
	{
		"Wow, it really works!",
		"AXQQEB10D5T20WK5C5P6RY90EXQQ4TVK44",
	},

	// from https://gist.github.com/szktty/228f85794e4187882a77734c89c384a8#gistcomment-3392026
	// > For example, both of `CR` and `CR0` can be decoded to `f`.
	{"f", "CR0"},

	// Clockwork-Base32 accepts both upper and lower case letters.
	{"foobar", "csqpyrk1e8"},
	{"Hello, world!", "91jprv3f5gg7evvjdhj22"},
	{
		"The quick brown fox jumps over the lazy dog.",
		"ahm6a83henmp6ts0c9s6yxve41k6yy10d9tptw3k41qqcsbj41t6gs90dhgqmy90chqpebg",
	},
	{
		"Wow, it really works!",
		"axqqeb10d5t20wk5c5p6ry90exqq4tvk44",
	},
	{"f", "cr"},
	{"f", "cr0"},
	{"fo", "csqg"},
	{"foo", "csqpy"},
	{"foob", "csqpyrg"},
	{"fooba", "csqpyrk1"},
	{"foobar", "csqpyrk1e8"},
	{
		"\x01\xdd\x3e\x62\xfe\x15\x4e\xd7\x2b\x6d\x2d\x24\x39\x74\x66\x9d",
		"07ekwrqy2n7deavd5mj3jx36km",
	},
	{
		"Wow, it really works!",
		"axqqeb10d5t20wk5c5p6ry90exqq4tvk44",
	},

	// "O" and "0" are treated as the same character.
	{
		"Wow, it really works!",
		"AXQQEB1OD5T2OWK5C5P6RY9OEXQQ4TVK44",
	},
	{
		"Wow, it really works!",
		"axqqeb1od5t2owk5c5p6ry9oexqq4tvk44",
	},

	// "I", "L" and "1" are treated as the same character.
	{"foobar", "CSQPYRKIE8"},
	{"foobar", "CSQPYRKLE8"},
	{"foobar", "csqpyrkie8"},
	{"foobar", "csqpyrkle8"},
}

func TestDecode(t *testing.T) {
	enc := NewEncoding()
	for _, testCase := range testCasesDecode {
		plain := make([]byte, enc.DecodedLen(len(testCase.encoded)))
		encoded := []byte(testCase.encoded)
		n, err := enc.Decode(plain, encoded)
		if err != nil {
			t.Errorf("error while decoding %q: %v", testCase.encoded, err)
		}
		if n != len(plain) {
			t.Errorf("unexpected length: want %d, got %d", len(plain), n)
		}
		if !bytes.Equal(plain, []byte(testCase.plain)) {
			t.Errorf("decoded %q, expected %q, actual %q\n",
				testCase.encoded, testCase.plain, plain)
		}
	}
}

var testCasesDecodeError = []struct {
	input string
	pos   int64
}{
	{"U", 0},
	{"u", 0},
	{"CSQG*", 4},
	{"CSQPYRK*", 7},
}

func TestDecode_Error(t *testing.T) {
	enc := NewEncoding()
	for _, testCase := range testCasesDecodeError {
		_, err := enc.DecodeString(testCase.input)
		switch err := err.(type) {
		case CorruptInputError:
			if int64(err) != testCase.pos {
				t.Errorf("unexpected error position: want %d, got %d", testCase.pos, int64(err))
			}
		default:
			t.Errorf("unexpected error type: want CorruptInputError, got %T", err)
		}
	}
}

func TestDecoder(t *testing.T) {
	enc := NewEncoding()
	for _, testCase := range testCasesDecode {
		r := strings.NewReader(testCase.encoded)
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(NewDecoder(enc, r)); err != nil {
			t.Errorf("error while decoding %q: %v", testCase.encoded, err)
		}
		plain := buf.Bytes()
		if !bytes.Equal(plain, []byte(testCase.plain)) {
			t.Errorf("decoded %q, expected %q, actual %q\n",
				testCase.encoded, testCase.plain, plain)
		}
	}
}

func TestDecoder_Buffering(t *testing.T) {
	for bs := 1; bs <= 24; bs++ {
		decoder := NewDecoder(Base32, strings.NewReader(bigtest.encoded))
		buf := make([]byte, len(bigtest.plain)+24)
		var total int
		var n int
		var err error
		for total = 0; total < len(bigtest.plain) && err == nil; {
			n, err = decoder.Read(buf[total : total+bs])
			total += n
		}
		if err != nil && err != io.EOF {
			t.Errorf("error while decoding %q: %v", bigtest.encoded, err)
		}
		if string(buf[0:total]) != bigtest.plain {
			t.Errorf("decoded %q, expected %q, actual %q\n",
				bigtest.encoded, bigtest.plain, buf[0:total])
		}
	}
}

func TestBig(t *testing.T) {
	n := 3*1000 + 1
	raw := make([]byte, n)
	const alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < n; i++ {
		raw[i] = alpha[i%len(alpha)]
	}
	encoded := new(bytes.Buffer)
	w := NewEncoder(Base32, encoded)
	nn, err := w.Write(raw)
	if nn != n || err != nil {
		t.Fatalf("Encoder.Write(raw) = %d, %v want %d, nil", nn, err, n)
	}
	err = w.Close()
	if err != nil {
		t.Fatalf("Encoder.Close() = %v want nil", err)
	}
	decoded, err := io.ReadAll(NewDecoder(Base32, encoded))
	if err != nil {
		t.Fatalf("io.ReadAll(NewDecoder(...)): %v", err)
	}

	if !bytes.Equal(raw, decoded) {
		var i int
		for i = 0; i < len(decoded) && i < len(raw); i++ {
			if decoded[i] != raw[i] {
				break
			}
		}
		t.Errorf("Decode(Encode(%d-byte string)) failed at offset %d", n, i)
	}
}

func BenchmarkEncode(b *testing.B) {
	data := make([]byte, 8192)
	buf := make([]byte, Base32.EncodedLen(len(data)))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		Base32.Encode(buf, data)
	}
}

func BenchmarkEncodeToString(b *testing.B) {
	data := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		Base32.EncodeToString(data)
	}
}

func BenchmarkDecode(b *testing.B) {
	data := make([]byte, Base32.EncodedLen(8192))
	Base32.Encode(data, make([]byte, 8192))
	buf := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		Base32.Decode(buf, data)
	}
}
func BenchmarkDecodeString(b *testing.B) {
	data := Base32.EncodeToString(make([]byte, 8192))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		Base32.DecodeString(data)
	}
}
