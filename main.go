package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/templexxx/xorsimd"
	"io"
	"log"
	"os"
)

var IV = []byte("3.141592653589793238462643383279")

func main() {
	if len(os.Args) == 0 {
		return
	}
	fi, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	out := IV    // the output of most recent invocation of the compression function; at the beginning this is the IV
	var b []byte // the currently used message block

	for {
		b = make([]byte, 32) // clear the value of b

		// read a chunk
		n, err := fi.Read(b)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// reached end of file; use the file size as the next block
		if n == 0 {
			s, err := fi.Stat()
			if err != nil {
				log.Fatal(err)
			}
			binary.PutVarint(b, s.Size())
		}

		// if less than 8 bytes were read add padding
		if n < 32 && n != 0 {
			// calculate how many bytes of padding are needed
			p := 32 - n

			// add padding to buf
			for i := 0; i < p; i++ {
				if i == 0 {
					b[n+i] = byte(1)
				} else {
					b[n+i] = byte(0)
				}
			}
		}

		out = f(out, b)

		if n == 0 {
			break
		}
	}
	fmt.Println(hex.EncodeToString(out))
}

// compression function f
func f(a []byte, b []byte) []byte {
	x := make([]byte, 32)
	xorsimd.Bytes(x, a, b)
	return x
}
