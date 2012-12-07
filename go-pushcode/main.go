package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
)

func copyToDropbox() (err error) {
	const dbp = "Dropbox/dev-go"
	for _, dropDirPath := range []string{filepath.Join("Q:", dbp), filepath.Join(ugo.UserHomeDirPath(), dbp)} {
		if uio.DirExists(dropDirPath) {
			if err = uio.ClearDirectory(dropDirPath); err == nil {
				if err = uio.CopyAll(ugo.GopathSrcGithub("metaleap"), filepath.Join(dropDirPath, "metaleap")); err == nil {
					err = uio.CopyAll(ugo.GopathSrcGithub("go3d"), filepath.Join(dropDirPath, "go3d"))
				}
			}
		}
	}
	return
}

func copyToRepos() (err error) {
	var (
		dirName, dirPath, srcDirPath string
		fileInfos                    []os.FileInfo
	)
	for _, repoBaseDirPath := range []string{"Q:\\gitrepos", "C:\\gitrepos"} {
		if fileInfos, _ = ioutil.ReadDir(repoBaseDirPath); len(fileInfos) > 0 {
			for _, fi := range fileInfos {
				if dirName = fi.Name(); fi.IsDir() {
					if srcDirPath = ugo.GopathSrcGithub("metaleap", dirName); !uio.DirExists(srcDirPath) {
						srcDirPath = ugo.GopathSrcGithub("go3d", dirName)
					}
					if dirPath = filepath.Join(repoBaseDirPath, dirName); uio.DirExists(srcDirPath) {
						if err = uio.ClearDirectory(dirPath, ".git"); err != nil {
							return
						} else if err = uio.CopyAll(srcDirPath, dirPath); err != nil {
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
	if err := copyToDropbox(); err != nil {
		panic(err)
	} else if err := copyToRepos(); err != nil {
		panic(err)
	}
}
