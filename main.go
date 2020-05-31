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

var IV = []byte("mdexpmnt")

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
		b = make([]byte, 8) // clear the value of b

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
			binary.LittleEndian.PutUint64(b, uint64(s.Size()))
		}

		// if less than 8 bytes were read add padding
		if n < 8 {
			// calculate how many bytes of padding are needed
			p := 8 - n

			// add padding to buf
			for i := 0; i < p; i++ {
				if i == 0 {
					b[n+i] = byte(1)
				} else {
					b[n+i] = byte(0)
				}
			}
		}

		out = f8(out, b)

		if n == 0 {
			break
		}
	}
	fmt.Println(hex.EncodeToString(out))
}

// compression function f for 8 byte blocks
func f8(a []byte, b []byte) []byte {
	x := make([]byte, 8)
	xorsimd.Bytes8(x, a, b)
	return x
}
