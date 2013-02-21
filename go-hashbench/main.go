package main

import (
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"runtime"
	"time"
)

var slice = make([]byte, 1)

func testHash(name string, h hash.Hash) {
	now := time.Now()
	for j := 0; j < 16; j++ {
		for slice[0] = 0; slice[0] < 255; slice[0]++ {
			h.Reset()
			h.Write(blob)
			h.Write(slice)
			h.Sum(nil)
		}
	}
	fmt.Printf("%s: %v\n", name, time.Now().Sub(now))
}

func main() {
	m := map[string]hash.Hash{
		"had":   adler32.New(),
		"hc32":  crc32.NewIEEE(),
		"hc64e": crc64.New(crc64.MakeTable(crc64.ECMA)),
		"hc64i": crc64.New(crc64.MakeTable(crc64.ISO)),
		"hf32":  fnv.New32(),
		"hf32a": fnv.New32a(),
		"hf64":  fnv.New64(),
		"hf64a": fnv.New64a(),
	}
	for n, h := range m {
		runtime.GC()
		testHash(n, h)
	}
}
