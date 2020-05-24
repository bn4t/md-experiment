package main

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/templexxx/xorsimd"
	"io/ioutil"
	"log"
	"os"
)

var IV = []byte("mdexpmnt")

func main() {
	if len(os.Args) == 0 {
		return
	}
	M, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	
	// get the length of M before the padding is added
	l := int64(len(M))
	
	// calculate how many bytes of padding are needed
	p := 8-len(M)%8
	
	// add padding to M
	for i:=0;i<p;i++{
		if i == 0 {
			M = append(M, byte(1))
		} else {
			M = append(M, byte(0))
		}
	}

	// get l as byte slice
	// See: https://stackoverflow.com/a/35371760
	lb := make([]byte, 8)
	binary.LittleEndian.PutUint64(lb, uint64(l))

	// append the length to the padded M
	I := append(M, []byte(lb)...)
	

	// prepare a slice to store all the message blocks
	var blocks [][]byte
	
	// split I into message blocks (b) of 8 bytes
	for b := I ;len(b)>=8; b = b[8:] {
		blocks = append(blocks, b[:8])
	}

	log.Print("padding (bytes): ",p)

	var out []byte
	
	// xor all the blocks
	// the first block is xored with the IV
	for k, v := range blocks {
		if k == 0 {
			xored := make([]byte, 8)
			xorsimd.Bytes8(xored, IV,v)
			out = xored
		} else {
			xored := make([]byte, 8)
			xorsimd.Bytes8(xored, out, v)
			out = xored
		}
	}
	log.Print(hex.EncodeToString(out))
}
