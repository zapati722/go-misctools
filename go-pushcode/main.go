package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
	ustr "github.com/metaleap/go-util/str"
)

var (
	pushToDropBox  = true
	pushToGitRepos = true

	wg            sync.WaitGroup
	dirTmpSkipper ustr.Matcher
)

func copyToDropbox() (err error) {
	const dbp = "Dropbox/dev-go"
	dropboxDirs := []string{"metaleap", "go3d", "goforks", "openbase", "go-leansites"}
	for _, dropDirPath := range []string{filepath.Join("Q:", dbp), filepath.Join(ugo.UserHomeDirPath(), dbp)} {
		if uio.DirExists(dropDirPath) {
			if err = uio.ClearDirectory(dropDirPath); err == nil {
				for _, githubName := range dropboxDirs {
					wg.Add(1)
					go func(dropDirPath, githubName string) {
						defer wg.Done()
						dst := filepath.Join(dropDirPath, githubName)
						if err := uio.CopyAll(ugo.GopathSrcGithub(githubName), dst, &dirTmpSkipper); err != nil {
							log.Printf("Error @ %s:\t%+v", dst, err)
						}
					}(dropDirPath, githubName)
				}
			}
			break
		}
	}
	return
}

func copyToRepos() (err error) {
	var (
		dirName, dirPath, srcDirPath string
		fileInfos                    []os.FileInfo
	)
	repoDirs := []string{"metaleap", "go3d", "goforks", "openbase"}
	for _, repoBaseDirPath := range []string{"Q:\\gitrepos", "C:\\gitrepos"} {
		if fileInfos, _ = ioutil.ReadDir(repoBaseDirPath); len(fileInfos) > 0 {
			for _, fi := range fileInfos {
				if dirName = fi.Name(); fi.IsDir() {
					for _, githubName := range repoDirs {
						if srcDirPath = ugo.GopathSrcGithub(githubName, dirName); uio.DirExists(srcDirPath) {
							break
						}
					}
					if dirPath = filepath.Join(repoBaseDirPath, dirName); uio.DirExists(srcDirPath) {
						wg.Add(1)
						go func(dirPath, srcDirPath string) {
							defer wg.Done()
							var err error
							if err = uio.ClearDirectory(dirPath, ".git"); err == nil {
								err = uio.CopyAll(srcDirPath, dirPath, &dirTmpSkipper)
							}
							if err != nil {
								log.Printf("Error @ %s:\t%+v", dirPath, err)
							}
						}(dirPath, srcDirPath)
					}
				}
			}
		}
	}
	return
}

func main() {
	var err error
	runtime.GOMAXPROCS(runtime.NumCPU())
	dirTmpSkipper.AddPatterns("_tmp")
	push := func(msg string, push bool, pusher func() error) {
		log.Println(msg)
		if push {
			if err = pusher(); err != nil {
				panic(err)
			} else {
				log.Println("YUP.")
			}
		} else {
			log.Println("NOPE.")
		}
	}
	push("DropBox?... ", pushToDropBox, copyToDropbox)
	push("GitHub?... ", pushToGitRepos, copyToRepos)
	log.Println("Waiting...")
	wg.Wait()
	log.Println("...all done.")
}
