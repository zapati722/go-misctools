package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
	ustr "github.com/metaleap/go-util/str"
)

var (
	pushToDropBox  = true
	pushToGitRepos = true

	dirTmpSkipper ustr.Matcher
)

func copyToDropbox() (err error) {
	const dbp = "Dropbox/dev-go"
	dropboxDirs := []string{"metaleap", "go3d", "openbase", "proxsite"}
	for _, dropDirPath := range []string{filepath.Join("Q:", dbp), filepath.Join(ugo.UserHomeDirPath(), dbp)} {
		if uio.DirExists(dropDirPath) {
			if err = uio.ClearDirectory(dropDirPath); err == nil {
				for _, githubName := range dropboxDirs {
					if err = uio.CopyAll(ugo.GopathSrcGithub(githubName), filepath.Join(dropDirPath, githubName), &dirTmpSkipper); err != nil {
						return
					}
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
	repoDirs := []string{"metaleap", "go3d", "openbase"}
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
						if err = uio.ClearDirectory(dirPath, ".git"); err != nil {
							return
						} else if err = uio.CopyAll(srcDirPath, dirPath, &dirTmpSkipper); err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}

func main() {
	var err error
	dirTmpSkipper.AddPatterns("_tmp")
	push := func(msg string, push bool, pusher func() error) {
		print(msg)
		if push {
			if err = pusher(); err != nil {
				panic(err)
			} else {
				println("YUP.")
			}
		} else {
			println("NOPE.")
		}
	}
	push("DropBox?... ", pushToDropBox, copyToDropbox)
	push("GitHub?... ", pushToGitRepos, copyToRepos)
}
