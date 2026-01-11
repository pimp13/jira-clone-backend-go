package main

import (
	"crypto/sha256"
	"log"
)

func main() {
	s := sha256.New()
	r, _ := s.Write([]byte("hefffllo"))
	e := s.Sum([]byte("dfdfd"))
	f := s.Size()
	log.Println(r, e, f)
}
