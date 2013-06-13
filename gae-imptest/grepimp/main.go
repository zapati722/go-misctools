package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
)

func main() {
	dir := flag.String("dir", filepath.Join(os.Getenv("GAESDK"), "goroot", "src", "pkg"), "")
	out := flag.String("out", ugo.GopathSrcGithub("metaleap", "go-misctools", "gae-imptest", "imp.txt"), "")
	flag.Parse()
	file, err := os.Create(*out)
	if err == nil {
		defer file.Close()
		m := map[string]bool{}
		if errs := uio.WalkAllFiles(*dir, func(fullPath string) bool {
			if filepath.Ext(fullPath) == ".go" {
				if relPath := filepath.Clean(filepath.Dir(fullPath)[len(*dir):]); !m[relPath] {
					m[relPath] = true
					if relPath = strings.Trim(filepath.ToSlash(relPath), "/"); relPath != "syscall" && relPath != "unsafe" && relPath != "builtin" && relPath != "runtime/race" && relPath != "net/http" && !(strings.Contains(relPath, "/testdata") || strings.Contains(relPath, "appengine_internal")) {
						file.WriteString(fmt.Sprintf("\t_ %#v\n", relPath))
					}
				}
			}
			return true
		}); len(errs) > 0 {
			err = errs[0]
		}
	}
	if err != nil {
		panic(err)
	}
}
