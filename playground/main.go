package main

import (
	"github.com/fastgh/go-ufs"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	f, _ := ufs.CreateLockFile(fs, "/tmp/hi.pid")
	f.Close()
}
