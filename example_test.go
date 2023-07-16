package clockwork_test

import (
	"fmt"
	"os"

	"github.com/shogo82148/go-clockwork-base32"
)

func Example() {
	msg := "Hello, 世界"
	encoded := clockwork.Base32.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := clockwork.Base32.DecodeString(encoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(string(decoded))
	// Output:
	// 91JPRV3F5GGE9E4PWYARR
	// Hello, 世界
}

func ExampleEncoding_EncodeToString() {
	data := []byte("any + old & data")
	str := clockwork.Base32.EncodeToString(data)
	fmt.Println(str)
	// Output:
	// C5Q7J81B41QPRS104RG68RBMC4
}

func ExampleEncoding_Encode() {
	data := []byte("Hello, world!")
	dst := make([]byte, clockwork.Base32.EncodedLen(len(data)))
	clockwork.Base32.Encode(dst, data)
	fmt.Println(string(dst))
	// Output:
	// 91JPRV3F5GG7EVVJDHJ22
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

func ExampleEncoding_Decode() {
	str := "91JPRV3F5GG7EVVJDHJ22"
	dst := make([]byte, clockwork.Base32.DecodedLen(len(str)))
	n, err := clockwork.Base32.Decode(dst, []byte(str))
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	dst = dst[:n]
	fmt.Printf("%q\n", dst)
	// Output:
	// "Hello, world!"
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
