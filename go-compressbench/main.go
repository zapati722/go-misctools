package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"time"

	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
)

type writerMaker func(w io.Writer) io.WriteCloser
type readerMaker func(r io.Reader) io.ReadCloser

var (
	blobs       [][]byte
	dirPath     = ugo.GopathSrcGithub("go3d", "go-ngine", "_examples", "-app-data", "_tmp", "tex")
	packStats   = map[string][]float64{}
	unpackStats = map[string][]float64{}
)

func testComp(name string, wm writerMaker, rm readerMaker) {
	var (
		err      error
		buf      *bytes.Buffer
		blob     []byte
		packed   [][]byte
		packer   io.WriteCloser
		unpacker io.ReadCloser
		now      time.Time
	)
	for _, blob = range blobs {
		runtime.GC()
		now = time.Now()
		buf = new(bytes.Buffer)
		packer = wm(buf)
		if _, err = packer.Write(blob); err != nil {
			panic(err)
		}
		if err = packer.Close(); err != nil {
			panic(err)
		}
		packStats[name] = append(packStats[name], time.Now().Sub(now).Seconds())
		packed = append(packed, buf.Bytes())
	}

	for _, blob = range packed {
		runtime.GC()
		now = time.Now()
		buf = bytes.NewBuffer(blob)
		unpacker = rm(buf)
		if _, err = ioutil.ReadAll(unpacker); err != nil {
			panic(err)
		}
		if err = unpacker.Close(); err != nil {
			panic(err)
		}
		unpackStats[name] = append(unpackStats[name], time.Now().Sub(now).Seconds())
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	uio.NewDirWalker(false, nil, func(fullPath string) bool {
		blobs = append(blobs, uio.ReadBinaryFile(fullPath, true))
		return true
	}).Walk(dirPath)

	testComp("flate1", func(w io.Writer) (wc io.WriteCloser) {
		var err error
		if wc, err = flate.NewWriter(w, 1); err != nil {
			panic(err)
		}
		return
	}, flate.NewReader)
	testComp("flate9", func(w io.Writer) (wc io.WriteCloser) {
		var err error
		if wc, err = flate.NewWriter(w, 9); err != nil {
			panic(err)
		}
		return
	}, flate.NewReader)
	testComp("lzw\t", func(w io.Writer) io.WriteCloser {
		return lzw.NewWriter(w, lzw.MSB, 8)
	}, func(r io.Reader) io.ReadCloser {
		return lzw.NewReader(r, lzw.MSB, 8)
	})
	testComp("zlib", func(w io.Writer) io.WriteCloser {
		return zlib.NewWriter(w)
	}, func(r io.Reader) (rc io.ReadCloser) {
		var err error
		if rc, err = zlib.NewReader(r); err != nil {
			panic(err)
		}
		return
	})
	testComp("gzip", func(w io.Writer) io.WriteCloser {
		return gzip.NewWriter(w)
	}, func(r io.Reader) (rc io.ReadCloser) {
		var err error
		if rc, err = gzip.NewReader(r); err != nil {
			panic(err)
		}
		return
	})
	printStats("PACK:", packStats)
	printStats("UNPACK:", unpackStats)
}

func printStats(desc string, m map[string][]float64) {
	var dur, min, max, avg float64
	fmt.Println(desc)
	for name, durs := range m {
		fmt.Printf("%s\t", name)
		avg, min, max = 0, math.MaxFloat64, -1
		for _, dur = range durs {
			if dur > max {
				max = dur
			}
			if dur < min {
				min = dur
			}
			avg += dur
		}
		avg = avg / float64(len(durs))
		fmt.Printf("avg:\t%v\tmin:\t%v\tmax:\t%v\n", avg, min, max)
	}
}
