package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/fastgh/go-ufs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_CopyFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.MkdirP(fs, "/Test_CopyFile_happy/c1")
	ufs.WriteTextP(fs, "/Test_CopyFile_happy/c1/src.txt", "hello")

	ufs.MkdirP(fs, "/Test_CopyFile_happy/c2")
	_, err := ufs.CopyFile(fs, "/Test_CopyFile_happy/c1/src.txt", "/Test_CopyFile_happy/c2/dest.txt")
	a.NoError(err)

	actual := ufs.ReadTextP(fs, "/Test_CopyFile_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func Test_CopyFile_SourceFileNotFound(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.MkdirP(fs, "/Test_CopyFile_SourceFileNotFound/c1")
	ufs.MkdirP(fs, "/Test_CopyFile_SourceFileNotFound/c2")

	a.Panicsf(func() {
		ufs.CopyFileP(fs, "/Test_CopyFile_SourceFileNotFound/c1/src.txt", "/Test_CopyFile_SourceFileNotFound/c2/dest.txt")
	}, "file not exists: /Test_CopyFile_SourceFileNotFound/c1/src.txt")
}

func Test_Rename_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.MkdirP(fs, "/Test_Rename_happy/c1")
	ufs.WriteTextP(fs, "/Test_Rename_happy/c1/src.txt", "hello")

	ufs.MkdirP(fs, "/Test_Rename_happy/c2")
	ufs.RenameP(fs, "/Test_Rename_happy/c1/src.txt", "/Test_Rename_happy/c2/dest.txt")

	actual := ufs.ReadTextP(fs, "/Test_Rename_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func Test_ReadLines_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.MkdirP(fs, "/Test_ReadAsLines_happy")
	ufs.WriteLinesP(fs, "/Test_ReadAsLines_happy/f.txt",
		"line 1",
		"line 2")

	actual := ufs.ReadLinesP(fs, "/Test_ReadAsLines_happy/f.txt")
	a.Equal([]string{
		"line 1",
		"line 2",
	}, actual)
}

func Test_ListSuffixed_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.MkdirP(fs, "Test_ListFilesWithExt_happy")
	ufs.MkdirP(fs, "Test_ListFilesWithExt_happy/.")
	ufs.MkdirP(fs, "Test_ListFilesWithExt_happy/..")
	ufs.MkdirP(fs, "Test_ListFilesWithExt_happy/d.hosts.txt")

	ufs.WriteTextP(fs, "Test_ListFilesWithExt_happy/1.hosts.txt", "1")
	ufs.WriteTextP(fs, "Test_ListFilesWithExt_happy/2.hosts.txt.not", "2")
	ufs.WriteTextP(fs, "Test_ListFilesWithExt_happy/3.hosts.not.text", "3")
	ufs.WriteTextP(fs, "Test_ListFilesWithExt_happy/4_empty.hosts.txt", "")
	ufs.WriteTextP(fs, "Test_ListFilesWithExt_happy/.hosts.txt", "5")

	a.Equal(map[string]string{
		"1": filepath.Join("Test_ListFilesWithExt_happy", "1.hosts.txt"),
	}, ufs.ListSuffixedP(fs, "Test_ListFilesWithExt_happy", ".hosts.txt", true))

	a.Equal(map[string]string{
		"1":       filepath.Join("Test_ListFilesWithExt_happy", "1.hosts.txt"),
		"4_empty": filepath.Join("Test_ListFilesWithExt_happy", "4_empty.hosts.txt"),
	}, ufs.ListSuffixedP(fs, "Test_ListFilesWithExt_happy", ".hosts.txt", false))
}

func Test_ExtractTitle_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("abc", ufs.ExtractTitle("/Test_ExtractTitle_happy/abc.xyz"))
	a.Equal("abc", ufs.ExtractTitle("/Test_ExtractTitle_happy/abc"))
	a.Equal("", ufs.ExtractTitle("/Test_ExtractTitle_happy/.xyz"))
}

func Test_FileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.False(ufs.FileExists(fs, "/f.txt"))
	ufs.WriteTextIfNotFoundP(fs, "/f.txt", "blah")
	a.True(ufs.FileExists(fs, "/f.txt"))

	ufs.MkdirP(fs, "/d")
	a.Panics(func() { ufs.FileExistsP(fs, "/d") }, "expect /d be file, but it is directory")
}

func Test_EnsureFileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { ufs.EnsureFileExistsP(fs, "/F.txt") }, "file not found: %s")
	ufs.WriteTextIfNotFoundP(fs, "/F.txt", "blah")
	ufs.EnsureFileExistsP(fs, "/F.txt")

	ufs.MkdirP(fs, "/D")
	a.Panics(func() { ufs.EnsureFileExistsP(fs, "/D") }, "expect /D be file, but it is directory")
}

func Test_DirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.False(ufs.DirExists(fs, "/d"))
	ufs.MkdirP(fs, "/d")
	a.True(ufs.DirExists(fs, "/d"))

	ufs.WriteTextIfNotFoundP(fs, "/f.txt", "blah")
	a.Panics(func() { ufs.DirExistsP(fs, "/f.txt") }, "expect /f.txt be directory, but it is file")
}

func Test_EnsureDirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { ufs.EnsureDirExistsP(fs, "/D") }, "directory not found: %s")
	ufs.MkdirP(fs, "/D")
	ufs.EnsureDirExistsP(fs, "/D")

	ufs.WriteTextIfNotFoundP(fs, "/f.txt", "blah")
	a.Panics(func() { ufs.EnsureDirExistsP(fs, "/f.txt") }, "expect /f.txt be directory, but it is file")
}

func Test_RemoveFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.RemoveFileP(fs, "/f.html")
	ufs.WriteTextIfNotFoundP(fs, "/f.html", "<html></html>")
	a.True(ufs.FileExists(fs, "/f.html"))
	ufs.RemoveFileP(fs, "/f.html")
	a.False(ufs.FileExists(fs, "/f.html"))

	ufs.MkdirP(fs, "/D")
	a.Panics(func() { ufs.RemoveFileP(fs, "/D") }, "expect /D be file, but it is directory")
	a.True(ufs.DirExists(fs, "/D"))
}

func Test_RemoveDir_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.RemoveDirP(fs, "/d")
	ufs.MkdirP(fs, "/d")
	ufs.RemoveDirP(fs, "/d")

	ufs.WriteTextIfNotFoundP(fs, "/f.html", "<html></html>")
	a.Panics(func() { ufs.RemoveDirP(fs, "/f.html") }, "expect /f.html be directory, but it is file")
}

func Test_WriteIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(ufs.WriteIfNotFoundP(fs, "/f.txt", []byte("hello")))
	a.Equal("hello", ufs.ReadTextP(fs, "/f.txt"))

	a.False(ufs.WriteIfNotFound(fs, "/f.txt", []byte("updated")))
	a.Equal("hello", ufs.ReadTextP(fs, "/f.txt"))
}

func Test_WriteTextIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(ufs.WriteTextIfNotFoundP(fs, "/f.txt", "hello"))
	a.Equal("hello", ufs.ReadTextP(fs, "/f.txt"))

	a.False(ufs.WriteTextIfNotFoundP(fs, "/f.txt", "updated"))
	a.Equal("hello", ufs.ReadTextP(fs, "/f.txt"))
}

func Test_TempFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	actual := ufs.TempFileP(fs, "xyz")
	a.True(strings.Contains(actual, "xyz"))
	a.NotEqual("xyz", actual)
}

func Test_ExpandHomePath_happy(t *testing.T) {
	a := require.New(t)
	a.Equal("none", ufs.ExpandHomePathP("none"))

	a.Equal(ufs.UserHomeDirP(), ufs.ExpandHomePathP("~"))
}
