# go-clockwork-base32

![Test](https://github.com/shogo82148/go-clockwork-base32/workflows/Test/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shogo82148/go-clockwork-base32)](https://pkg.go.dev/github.com/shogo82148/go-clockwork-base32)

An Implementation Clockwork-Base32 for Go.

Clockwork Base32 is a simple variant of Base32 inspired by Crockford's Base32.
See [Clockwork Base32 Specification](https://gist.github.com/szktty/228f85794e4187882a77734c89c384a8).

## Usage

The interface is compatible with [encoding/base32](https://golang.org/pkg/encoding/base32/).

```go
func ExampleEncoding_EncodeToString() {
	data := []byte("any + old & data")
	str := clockwork.Base32.EncodeToString(data)
	fmt.Println(str)
	// Output:
	// C5Q7J81B41QPRS104RG68RBMC4
}

func ExampleEncoding_DecodeString() {
	str := "EDQPTS90CHGQ8R90EXMQ8T1000G62VK443QVQFR"
	data, err := clockwork.Base32.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%q\n", data)
	// Output:
	// "some data with \x00 and \ufeff"
}

func ExampleNewEncoder() {
	input := []byte("foo\x00bar")
	encoder := clockwork.NewEncoder(clockwork.Base32, os.Stdout)
	encoder.Write(input)
	// Must close the encoder when finished to flush any partial blocks.
	// If you comment out the following line, the last partial block "r"
	// won't be encoded.
	encoder.Close()
	// Output:
	// CSQPY032C5S0
}
```

## See Also

- [Clockwork Base32 Specification](https://gist.github.com/szktty/228f85794e4187882a77734c89c384a8)
- [szktty/go-clockwork-base32](https://github.com/szktty/go-clockwork-base32)
    - A reference implementation of Clockwork Base32 for Go.
- [encoding/base32](https://golang.org/pkg/encoding/base32/)
    - Go standard library of RFC 4648 base32 encoding
