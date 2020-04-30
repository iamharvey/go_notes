package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	h := sha1.New()
	b1 := []byte{0x00}
	b2 := []byte("hello sha1\n")
	b1 = append(b1, b2...)
	h.Write(b1)
	s := h.Sum(nil)
	fmt.Printf("%x\n", s)
}

