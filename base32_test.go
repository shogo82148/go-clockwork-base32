package clockwork

import (
	"bytes"
	"testing"
)

type testCase struct {
	plain   string
	encoded string
}

var testCases = []testCase{
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
}

func TestEncode(t *testing.T) {
	enc := NewEncoding()
	for _, testCase := range testCases {
		plain := []byte(testCase.plain)
		encoded := make([]byte, enc.EncodedLen(len(plain)))
		enc.Encode(encoded, plain)
		if bytes.Compare(encoded, []byte(testCase.encoded)) != 0 {
			t.Errorf("encoded '%s', expected '%s', actual '%s'\n",
				testCase.plain, testCase.encoded, encoded)
		}
	}
}
