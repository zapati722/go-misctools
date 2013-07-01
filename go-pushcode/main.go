package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-utils/ufs"
	"github.com/go-utils/ugo"
	"github.com/go-utils/ustr"
)

var (
	pushToDropBox  = true
	pushToGitRepos = true

	wg            sync.WaitGroup
	dirTmpSkipper ustr.Matcher
)

func copyToDropbox() (err error) {
	const dbp = "Dropbox/dev-go"
	dropboxDirs := []string{"metaleap", "go3d", "go-forks", "go-utils", "openbase", "go-leansites"}
	for _, dropDirPath := range []string{filepath.Join("Q:", dbp), filepath.Join(ugo.UserHomeDirPath(), dbp)} {
		if ufs.DirExists(dropDirPath) {
			if err = ufs.ClearDirectory(dropDirPath); err == nil {
				for _, githubName := range dropboxDirs {
					wg.Add(1)
					go func(dropDirPath, githubName string) {
						defer wg.Done()
						dst := filepath.Join(dropDirPath, githubName)
						if err := ufs.CopyAll(ugo.GopathSrcGithub(githubName), dst, &dirTmpSkipper); err != nil {
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
	repoDirs := []string{"metaleap", "go3d", "openbase", "go-utils"}
	for _, repoBaseDirPath := range []string{"Q:\\gitrepos", "C:\\gitrepos"} {
		if fileInfos, _ = ioutil.ReadDir(repoBaseDirPath); len(fileInfos) > 0 {
			for _, fi := range fileInfos {
				if dirName = fi.Name(); fi.IsDir() {
					for _, githubName := range repoDirs {
						if srcDirPath = ugo.GopathSrcGithub(githubName, dirName); ufs.DirExists(srcDirPath) {
							break
						}
					}
					if dirPath = filepath.Join(repoBaseDirPath, dirName); ufs.DirExists(srcDirPath) {
						wg.Add(1)
						go func(dirPath, srcDirPath string) {
							defer wg.Done()
							var err error
							if err = ufs.ClearDirectory(dirPath, ".git"); err == nil {
								err = ufs.CopyAll(srcDirPath, dirPath, &dirTmpSkipper)
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
	ugo.MaxProcs()
	dirTmpSkipper.AddPatterns("_tmp")
	push := func(msg string, push bool, pusher func() error) {
		defer wg.Done()
		if push {
			if err = pusher(); err != nil {
				panic(err)
			} else {
				log.Println(msg + "YUP.")
			}
		} else {
			log.Println(msg + "NOPE.")
		}
	}
	wg.Add(2)
	log.Println("Waiting...")
	go push("DropBox?... ", pushToDropBox, copyToDropbox)
	go push("GitHub?... ", pushToGitRepos, copyToRepos)
	wg.Wait()
	log.Println("...all done.")
}
