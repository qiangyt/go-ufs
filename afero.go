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

	"github.com/fastgh/go-comm/v2"
	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

func DefaultEtcHostsP() string {
	r, err := DefaultEtcHosts()
	if err != nil {
		panic(err)
	}
	return r
}

func DefaultEtcHosts() (string, error) {
	if comm.IsWindows() {
		return `C:\Windows\System32\Drivers\etc\hosts`, nil
	}
	if comm.IsLinux() {
		return "/etc/hosts", nil
	}
	if comm.IsDarwin() {
		return "/private/etc/hosts", nil
	}
	return "", errors.New(runtime.GOOS + " is not yet supported")
}

func CopyFileP(fs afero.Fs, path string, newPath string) int64 {
	r, err := CopyFile(fs, path, newPath)
	if err != nil {
		panic(err)
	}
	return r
}

func CopyFile(fs afero.Fs, path string, newPath string) (int64, error) {
	EnsureFileExists(fs, path)

	src, err := fs.Open(path)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to read file %s", path)
	}
	defer src.Close()

	dst, err := fs.Create(newPath)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create file %s", newPath)
	}
	defer dst.Close()

	nBytes, err := io.Copy(dst, src)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to copy file %s to %s", path, newPath)
	}
	return nBytes, nil
}

func RenameP(fs afero.Fs, path string, newPath string) {
	if err := Rename(fs, path, newPath); err != nil {
		panic(err)
	}
}

func Rename(fs afero.Fs, path string, newPath string) error {
	err := fs.Rename(path, newPath)
	if err != nil {
		return errors.Wrapf(err, "failed to move file %s to %s", path, newPath)
	}
	return nil
}

func StatP(fs afero.Fs, path string, ensureExists bool) os.FileInfo {
	r, err := Stat(fs, path, ensureExists)
	if err != nil {
		panic(err)
	}
	return r
}

// Stat ...
func Stat(fs afero.Fs, path string, ensureExists bool) (os.FileInfo, error) {
	r, err := fs.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "failed to stat file: %s", path)
		}
		if ensureExists {
			return nil, errors.Wrapf(err, "file not found: %s", path)
		}
		return nil, nil
	}

	return r, nil
}

func FileExistsP(fs afero.Fs, path string) bool {
	r, err := FileExists(fs, path)
	if err != nil {
		panic(err)
	}
	return r
}

// FileExists ...
func FileExists(fs afero.Fs, path string) (bool, error) {
	fi, err := Stat(fs, path, false)
	if err != nil {
		return false, err
	}
	if fi == nil {
		return false, nil
	}
	if fi.IsDir() {
		return false, fmt.Errorf("expect %s be file, but it is directory", path)
	}
	return true, nil
}

func EnsureFileExistsP(fs afero.Fs, path string) {
	if err := EnsureFileExists(fs, path); err != nil {
		panic(err)
	}
}

func EnsureFileExists(fs afero.Fs, path string) error {
	exists, err := FileExists(fs, path)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file not found: %s", path)
	}
	return nil
}

func DirExistsP(fs afero.Fs, path string) bool {
	r, err := DirExists(fs, path)
	if err != nil {
		panic(err)
	}
	return r
}

// DirExists ...
func DirExists(fs afero.Fs, path string) (bool, error) {
	fi, err := Stat(fs, path, false)
	if err != nil {
		return false, err
	}
	if fi == nil {
		return false, nil
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("expect %s be directory, but it is file", path)
	}
	return true, nil
}

func EnsureDirExistsP(fs afero.Fs, path string) {
	if err := EnsureDirExists(fs, path); err != nil {
		panic(err)
	}
}

func EnsureDirExists(fs afero.Fs, path string) error {
	exists, err := DirExists(fs, path)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("directory not found: %s", path)
	}
	return nil
}

func RemoveFileP(fs afero.Fs, path string) {
	if err := RemoveFile(fs, path); err != nil {
		panic(err)
	}
}

// RemoveFile ...
func RemoveFile(fs afero.Fs, path string) error {
	found, err := FileExists(fs, path)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	if err := fs.Remove(path); err != nil {
		return errors.Wrapf(err, "failed to delete file: %s", path)
	}
	return nil
}

func RemoveDirP(fs afero.Fs, path string) {
	if err := RemoveDir(fs, path); err != nil {
		panic(err)
	}
}

// RemoveDir ...
func RemoveDir(fs afero.Fs, path string) error {
	if path == "/" || path == "\\" {
		return fmt.Errorf("should NOT remove root directory")
	}
	found, err := DirExists(fs, path)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	if err := fs.RemoveAll(path); err != nil {
		return errors.Wrapf(err, "failed to delete directory: %s", path)
	}
	return nil
}

func ReadBytesP(fs afero.Fs, path string) []byte {
	r, err := ReadBytes(fs, path)
	if err != nil {
		return r
	}
	return r
}

// ReadBytes ...
func ReadBytes(fs afero.Fs, path string) ([]byte, error) {
	r, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", path)
	}
	return r, nil
}

func ReadTextP(fs afero.Fs, path string) string {
	r, err := ReadText(fs, path)
	if err != nil {
		panic(err)
	}
	return r
}

func ReadText(fs afero.Fs, path string) (string, error) {
	bytes, err := ReadBytes(fs, path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ReadLinesP(fs afero.Fs, path string) []string {
	r, err := ReadLines(fs, path)
	if err != nil {
		panic(err)
	}
	return r
}

func ReadLines(fs afero.Fs, path string) ([]string, error) {
	if err := EnsureFileExists(fs, path); err != nil {
		return nil, err
	}

	f, err := fs.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	return comm.ReadLines(f), nil
}

func WriteIfNotFoundP(fs afero.Fs, path string, content []byte) bool {
	r, err := WriteIfNotFound(fs, path, content)
	if err != nil {
		panic(err)
	}
	return r
}

// WriteIfNotFound ...
func WriteIfNotFound(fs afero.Fs, path string, content []byte) (bool, error) {
	found, err := FileExists(fs, path)
	if err != nil {
		return false, err
	}
	if found {
		return false, nil
	}
	if err := Write(fs, path, content); err != nil {
		return true, err
	}
	return true, nil
}

func WriteP(fs afero.Fs, path string, content []byte) {
	if err := Write(fs, path, content); err != nil {
		panic(err)
	}
}

// Write ...
func Write(fs afero.Fs, path string, content []byte) error {
	if err := afero.WriteFile(fs, path, content, 0o640); err != nil {
		return errors.Wrapf(err, "failed to write file: %s", path)
	}
	return nil
}

func WriteTextP(fs afero.Fs, path string, content string) {
	if err := WriteText(fs, path, content); err != nil {
		panic(err)
	}
}

// WriteText ...
func WriteText(fs afero.Fs, path string, content string) error {
	return Write(fs, path, []byte(content))
}

// WriteTextIfNotFound ...
func WriteTextIfNotFoundP(fs afero.Fs, path string, content string) bool {
	r, err := WriteTextIfNotFound(fs, path, content)
	if err != nil {
		panic(err)
	}
	return r
}

func WriteTextIfNotFound(fs afero.Fs, path string, content string) (bool, error) {
	found, err := FileExists(fs, path)
	if err != nil {
		return found, err
	}
	if found {
		return false, nil
	}
	if err := WriteText(fs, path, content); err != nil {
		return false, err
	}
	return true, nil
}

// WriteLines ...
func WriteLinesP(fs afero.Fs, path string, lines ...string) {
	if err := WriteLines(fs, path, lines...); err != nil {
		panic(err)
	}
}

func WriteLines(fs afero.Fs, path string, lines ...string) error {
	return WriteText(fs, path, comm.JoinedLines(lines...))
}

func ExpandHomePathP(path string) string {
	r, err := ExpandHomePath(path)
	if err != nil {
		panic(err)
	}
	return r
}

// ExpandHomePath ...
func ExpandHomePath(path string) (string, error) {
	var r string
	var err error

	if r, err = homedir.Expand(path); err != nil {
		return "", errors.Wrapf(err, "failed to expand path: %s", path)
	}
	return r, nil
}

func UserHomeDirP() string {

	r, err := UserHomeDir()
	if err != nil {
		panic(err)
	}
	return r

}

func UserHomeDir() (string, error) {
	r, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home path")
	}
	return r, nil
}

func MkdirP(fs afero.Fs, path string) {
	if err := Mkdir(fs, path); err != nil {
		panic(err)
	}
}

// Mkdir ...
func Mkdir(fs afero.Fs, path string) error {
	if err := fs.MkdirAll(path, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create directory: %s", path)
	}
	return nil
}

func ListSuffixedP(fs afero.Fs, targetDir string, suffix string, skipEmptyFile bool) map[string]string {
	r, err := ListSuffixed(fs, targetDir, suffix, skipEmptyFile)
	if err != nil {
		panic(err)
	}
	return r
}

func ListSuffixed(fs afero.Fs, targetDir string, suffix string, skipEmptyFile bool) (map[string]string, error) {
	fiList, err := afero.ReadDir(fs, targetDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read directory: %s", targetDir)
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

	return r, nil
}

func ExtractTitle(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	return base[:len(base)-len(ext)]
}

func TempFileP(fs afero.Fs, pattern string) string {
	r, err := TempFile(fs, pattern)
	if err != nil {
		panic(err)
	}
	return r
}

func TempFile(fs afero.Fs, pattern string) (string, error) {
	f, err := afero.TempFile(fs, "", pattern)
	if err != nil {
		return "", errors.Wrap(err, "failed to create temporary file")
	}
	r := f.Name()
	f.Close()

	return r, nil
}

func TempTextFileP(fs afero.Fs, pattern string, content string) string {
	r, err := TempTextFile(fs, pattern, content)
	if err != nil {
		panic(err)
	}
	return r
}

func TempTextFile(fs afero.Fs, pattern string, content string) (string, error) {
	r, err := TempFile(fs, pattern)
	if err != nil {
		return "", err
	}
	if err := WriteText(fs, r, content); err != nil {
		return "", err
	}
	return r, nil
}
