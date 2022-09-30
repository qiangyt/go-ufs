//go:build windows
// +build windows

package ufs

// a revised copy of github.com/allan-simon/go-singleinstance v0.0.0-20210120080615-d0997106ab37

import (
	"strconv"

	"github.com/spf13/afero"
)

// CreateLockFile tries to create a file with given name and acquire an
// exclusive lock on it. If the file already exists AND is still locked, it will
// fail.
func CreateLockFile(fs afero.Fs, filename string) (afero.File, error) {
	if _, err := fs.Stat(filename); err == nil {
		// If the files exists, we first try to remove it
		if err = fs.Remove(filename); err != nil {
			return nil, err
		}
	} else if !fs.IsNotExist(err) {
		return nil, err
	}

	f, err := fs.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	// Write PID to lock file
	_, err = f.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		return nil, err
	}

	return f, nil
}
