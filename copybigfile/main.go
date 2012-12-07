package main

import (
	"flag"
	"io"
	"log"
	"os"
)

var (
	flagSrcFilePath = flag.String("src", "", "source file path")
	flagDstFilePath = flag.String("dst", "", "destination file path")
)

func main() {
	flag.Parse()
	src, err := os.Open(*flagSrcFilePath)
	if err != nil {
		panic(err)
	}
	defer src.Close()
	dst, err := os.Create(*flagDstFilePath)
	if err != nil {
		panic(err)
	}
	defer dst.Close()
	eof, errs, n, nn, l := false, 0, 0, 0, int64(1*1024*1024*1024)
	i, off, buf, one := int64(0), int64(0), make([]byte, l), make([]byte, 1)
	for j := 0; !eof; j++ {
		log.Printf("%vGB done. Reading 1GB at %v...", j, off)
		n, err = src.ReadAt(buf, off)
		if eof = (err == io.EOF); (err != nil) && !eof {
			log.Printf("ERR: %v; now reading this 1GB chunk byte-for-byte...", err)
			n, errs = 0, 0
			for i = 0; i < l; i++ {
				if nn, err = src.ReadAt(one, i+off); nn == 1 {
					buf[i] = one[0]
				} else if eof = (err == io.EOF); eof {
					break
				} else {
					errs++
					buf[i] = 0
				}
				n++
			}
			log.Printf("Done byte-for-byte reading of this 1GB chunk, used the value 0 for %v corrupted / non-readable bytes.", errs)
		}
		if n <= 0 {
			eof = true
		} else {
			log.Printf("Writing %v bytes...", n)
			if _, err = dst.Write(buf[:n]); err != nil {
				panic(err)
			}
			off += int64(n)
		}
	}
}
