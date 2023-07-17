package clockwork

import "testing"

func FuzzEncode(f *testing.F) {
	for _, t := range testCasesEncode {
		f.Add(t.plain)
	}
	f.Fuzz(func(t *testing.T, input string) {
		enc := NewEncoding()
		encoded := enc.EncodeToString([]byte(input))
		decoded, err := enc.DecodeString(encoded)
		if err != nil {
			t.Error(err)
		}
		if string(decoded) != input {
			t.Errorf("decoded string does not match input: %q != %q", string(decoded), input)
		}
	})
}

func FuzzDecode(f *testing.F) {
	for _, t := range testCasesDecode {
		f.Add(t.encoded)
	}
	f.Fuzz(func(t *testing.T, a string) {
		enc := NewEncoding()
		_, err := enc.DecodeString(a)
		if err != nil {
			return
		}
	})
}
