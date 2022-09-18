package ufs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"

	"github.com/fastgh/go-comm"
	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

func DefaultEtcHosts() string {
	switch comm.DefaultOSType() {
	case comm.Windows:
		return `C:\Windows\System32\Drivers\etc\hosts`
	case comm.Linux:
		return "/etc/hosts"
	case comm.Darwin:
		return "/private/etc/hosts"
	default:
		panic(errors.New(runtime.GOOS + " is not yet supported"))
	}
}

func CopyFile(fs afero.Fs, path string, newPath string) int64 {
	EnsureFileExists(fs, path)

	src, err := fs.Open(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read file %s", path))
	}
	defer src.Close()

	dst, err := fs.Create(newPath)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create file %s", newPath))
	}
	defer dst.Close()

	nBytes, err := io.Copy(dst, src)
	if err != nil {
		panic(errors.Wrapf(err, "failed to copy file %s to %s", path, newPath))
	}
	return nBytes
}

func Rename(fs afero.Fs, path string, newPath string) {
	err := fs.Rename(path, newPath)
	if err != nil {
		panic(errors.Wrapf(err, "failed to move file %s to %s", path, newPath))
	}
}

// Stat ...
func Stat(fs afero.Fs, path string, ensureExists bool) os.FileInfo {
	r, err := fs.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(errors.Wrapf(err, "failed to stat file: %s", path))
		}
		if ensureExists {
			panic(errors.Wrapf(err, "file not exists: %s", path))
		}
		return nil
	}

	return r
}

// FileExists ...
func FileExists(fs afero.Fs, path string) bool {
	fi := Stat(fs, path, false)
	if fi == nil {
		return false
	}
	if fi.IsDir() {
		panic(fmt.Errorf("expect %s be file, but it is directory", path))
	}
	return true
}

func EnsureFileExists(fs afero.Fs, path string) {
	if !FileExists(fs, path) {
		panic(fmt.Errorf("file not found: %s", path))
	}
}

// DirExists ...
func DirExists(fs afero.Fs, path string) bool {
	fi := Stat(fs, path, false)
	if fi == nil {
		return false
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("expect %s be directory, but it is file", path))
	}
	return true
}

func EnsureDirExists(fs afero.Fs, path string) {
	if !DirExists(fs, path) {
		panic(fmt.Errorf("directory not found: %s", path))
	}
}

// RemoveFile ...
func RemoveFile(fs afero.Fs, path string) {
	if FileExists(fs, path) {
		if err := fs.Remove(path); err != nil {
			panic(errors.Wrapf(err, "failed to delete file: %s", path))
		}
	}
}

// RemoveDir ...
func RemoveDir(fs afero.Fs, path string) {
	if path == "/" || path == "\\" {
		panic(fmt.Errorf("should NOT remove root directory"))
	}
	if DirExists(fs, path) {
		if err := fs.RemoveAll(path); err != nil {
			panic(errors.Wrapf(err, "failed to delete directory: %s", path))
		}
	}
}

// ReadBytes ...
func ReadBytes(fs afero.Fs, path string) []byte {
	r, err := afero.ReadFile(fs, path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read file: %s", path))
	}
	return r
}

func ReadText(fs afero.Fs, path string) string {
	bytes := ReadBytes(fs, path)
	return string(bytes)
}

func ReadLines(fs afero.Fs, path string) []string {
	EnsureFileExists(fs, path)

	f, err := fs.Open(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open file: %s", path))
	}
	defer f.Close()

	return comm.ReadLines(f)
}

// WriteIfNotFound ...
func WriteIfNotFound(fs afero.Fs, path string, content []byte) bool {
	if FileExists(fs, path) {
		return false
	}
	Write(fs, path, content)
	return true
}

// Write ...
func Write(fs afero.Fs, path string, content []byte) {
	if err := afero.WriteFile(fs, path, content, 0640); err != nil {
		panic(errors.Wrapf(err, "failed to write file: %s", path))
	}
}

// WriteText ...
func WriteText(fs afero.Fs, path string, content string) {
	Write(fs, path, []byte(content))
}

// WriteTextIfNotFound ...
func WriteTextIfNotFound(fs afero.Fs, path string, content string) bool {
	if FileExists(fs, path) {
		return false
	}
	WriteText(fs, path, content)
	return true
}

// WriteLines ...
func WriteLines(fs afero.Fs, path string, lines ...string) {
	WriteText(fs, path, comm.JoinedLines(lines...))
}

// WriteLinesIfNotFound ...
func WriteLinesIfNotFound(fs afero.Fs, path string, lines ...string) bool {
	if FileExists(fs, path) {
		return false
	}
	WriteLines(fs, path, lines...)
	return true
}

// ExpandHomePath ...
func ExpandHomePath(path string) string {
	var r string
	var err error

	if r, err = homedir.Expand(path); err != nil {
		panic(errors.Wrapf(err, "failed to expand path: %s", path))
	}
	return r
}

func UserHomeDir() string {
	r, err := os.UserHomeDir()
	if err != nil {
		panic(errors.Wrap(err, "failed to get home path"))
	}
	return r
}

// Mkdir ...
func Mkdir(fs afero.Fs, path string) {
	if err := fs.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "failed to create directory: %s", path))
	}
}

func ListSuffixed(fs afero.Fs, targetDir string, suffix string, skipEmptyFile bool) map[string]string {
	fiList, err := afero.ReadDir(fs, targetDir)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read directory: %s", targetDir))
	}

	extLen := len(suffix)

	r := map[string]string{}
	for _, fi := range fiList {
		if fi.IsDir() {
			continue
		}
		if skipEmptyFile && fi.Size() == 0 {
			continue
		}

		fBase := filepath.Base(fi.Name())
		if !strings.HasSuffix(fBase, suffix) {
			continue
		}
		if len(fBase) == extLen {
			continue
		}

		fTitle := fBase[:len(fBase)-extLen]
		r[fTitle] = filepath.Join(targetDir, fi.Name())
	}

	return r
}

func ExtractTitle(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	return base[:len(base)-len(ext)]
}

func TempFile(fs afero.Fs, pattern string) string {
	f, err := afero.TempFile(fs, "", pattern)
	if err != nil {
		panic(errors.Wrap(err, "failed to create temporary file"))
	}
	r := f.Name()
	f.Close()

	return r
}

func TempTextFile(fs afero.Fs, pattern string, content string) string {
	r := TempFile(fs, pattern)
	WriteText(fs, r, content)
	return r
}
