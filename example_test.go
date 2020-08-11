package clockwork_test

import (
	"fmt"
	"os"

	"github.com/shogo82148/go-clockwork-base32"
)

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
