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

	ufs.Mkdir(fs, "/Test_CopyFile_happy/c1")
	ufs.WriteText(fs, "/Test_CopyFile_happy/c1/src.txt", "hello")

	ufs.Mkdir(fs, "/Test_CopyFile_happy/c2")
	ufs.CopyFile(fs, "/Test_CopyFile_happy/c1/src.txt", "/Test_CopyFile_happy/c2/dest.txt")

	actual := ufs.ReadText(fs, "/Test_CopyFile_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func Test_CopyFile_SourceFileNotFound(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.Mkdir(fs, "/Test_CopyFile_SourceFileNotFound/c1")
	ufs.Mkdir(fs, "/Test_CopyFile_SourceFileNotFound/c2")

	a.Panicsf(func() {
		ufs.CopyFile(fs, "/Test_CopyFile_SourceFileNotFound/c1/src.txt", "/Test_CopyFile_SourceFileNotFound/c2/dest.txt")
	}, "file not exists: /Test_CopyFile_SourceFileNotFound/c1/src.txt")
}

func Test_Rename_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.Mkdir(fs, "/Test_Rename_happy/c1")
	ufs.WriteText(fs, "/Test_Rename_happy/c1/src.txt", "hello")

	ufs.Mkdir(fs, "/Test_Rename_happy/c2")
	ufs.Rename(fs, "/Test_Rename_happy/c1/src.txt", "/Test_Rename_happy/c2/dest.txt")

	actual := ufs.ReadText(fs, "/Test_Rename_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func Test_ReadLines_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.Mkdir(fs, "/Test_ReadAsLines_happy")
	ufs.WriteLines(fs, "/Test_ReadAsLines_happy/f.txt",
		"line 1",
		"line 2")

	actual := ufs.ReadLines(fs, "/Test_ReadAsLines_happy/f.txt")
	a.Equal([]string{
		"line 1",
		"line 2",
	}, actual)
}

func Test_ListSuffixed_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.Mkdir(fs, "Test_ListFilesWithExt_happy")
	ufs.Mkdir(fs, "Test_ListFilesWithExt_happy/.")
	ufs.Mkdir(fs, "Test_ListFilesWithExt_happy/..")
	ufs.Mkdir(fs, "Test_ListFilesWithExt_happy/d.hosts.txt")

	ufs.WriteText(fs, "Test_ListFilesWithExt_happy/1.hosts.txt", "1")
	ufs.WriteText(fs, "Test_ListFilesWithExt_happy/2.hosts.txt.not", "2")
	ufs.WriteText(fs, "Test_ListFilesWithExt_happy/3.hosts.not.text", "3")
	ufs.WriteText(fs, "Test_ListFilesWithExt_happy/4_empty.hosts.txt", "")
	ufs.WriteText(fs, "Test_ListFilesWithExt_happy/.hosts.txt", "5")

	a.Equal(map[string]string{
		"1": filepath.Join("Test_ListFilesWithExt_happy", "1.hosts.txt"),
	}, ufs.ListSuffixed(fs, "Test_ListFilesWithExt_happy", ".hosts.txt", true))

	a.Equal(map[string]string{
		"1":       filepath.Join("Test_ListFilesWithExt_happy", "1.hosts.txt"),
		"4_empty": filepath.Join("Test_ListFilesWithExt_happy", "4_empty.hosts.txt"),
	}, ufs.ListSuffixed(fs, "Test_ListFilesWithExt_happy", ".hosts.txt", false))
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
	ufs.WriteTextIfNotFound(fs, "/f.txt", "blah")
	a.True(ufs.FileExists(fs, "/f.txt"))

	ufs.Mkdir(fs, "/d")
	a.Panics(func() { ufs.FileExists(fs, "/d") }, "expect /d be file, but it is directory")
}

func Test_EnsureFileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { ufs.EnsureFileExists(fs, "/F.txt") }, "file not found: %s")
	ufs.WriteTextIfNotFound(fs, "/F.txt", "blah")
	ufs.EnsureFileExists(fs, "/F.txt")

	ufs.Mkdir(fs, "/D")
	a.Panics(func() { ufs.EnsureFileExists(fs, "/D") }, "expect /D be file, but it is directory")
}

func Test_DirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.False(ufs.DirExists(fs, "/d"))
	ufs.Mkdir(fs, "/d")
	a.True(ufs.DirExists(fs, "/d"))

	ufs.WriteTextIfNotFound(fs, "/f.txt", "blah")
	a.Panics(func() { ufs.DirExists(fs, "/f.txt") }, "expect /f.txt be directory, but it is file")
}

func Test_EnsureDirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { ufs.EnsureDirExists(fs, "/D") }, "directory not found: %s")
	ufs.Mkdir(fs, "/D")
	ufs.EnsureDirExists(fs, "/D")

	ufs.WriteTextIfNotFound(fs, "/f.txt", "blah")
	a.Panics(func() { ufs.EnsureDirExists(fs, "/f.txt") }, "expect /f.txt be directory, but it is file")
}

func Test_RemoveFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.RemoveFile(fs, "/f.html")
	ufs.WriteTextIfNotFound(fs, "/f.html", "<html></html>")
	a.True(ufs.FileExists(fs, "/f.html"))
	ufs.RemoveFile(fs, "/f.html")
	a.False(ufs.FileExists(fs, "/f.html"))

	ufs.Mkdir(fs, "/D")
	a.Panics(func() { ufs.RemoveFile(fs, "/D") }, "expect /D be file, but it is directory")
	a.True(ufs.DirExists(fs, "/D"))
}

func Test_RemoveDir_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.RemoveDir(fs, "/d")
	ufs.Mkdir(fs, "/d")
	ufs.RemoveDir(fs, "/d")

	ufs.WriteTextIfNotFound(fs, "/f.html", "<html></html>")
	a.Panics(func() { ufs.RemoveDir(fs, "/f.html") }, "expect /f.html be directory, but it is file")
}

func Test_WriteIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(ufs.WriteIfNotFound(fs, "/f.txt", []byte("hello")))
	a.Equal("hello", ufs.ReadText(fs, "/f.txt"))

	a.False(ufs.WriteIfNotFound(fs, "/f.txt", []byte("updated")))
	a.Equal("hello", ufs.ReadText(fs, "/f.txt"))
}

func Test_WriteTextIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(ufs.WriteTextIfNotFound(fs, "/f.txt", "hello"))
	a.Equal("hello", ufs.ReadText(fs, "/f.txt"))

	a.False(ufs.WriteTextIfNotFound(fs, "/f.txt", "updated"))
	a.Equal("hello", ufs.ReadText(fs, "/f.txt"))
}

func Test_WriteLinesIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(ufs.WriteLinesIfNotFound(fs, "/f.txt", "1", "2"))
	a.Equal("1\n2", ufs.ReadText(fs, "/f.txt"))

	a.False(ufs.WriteLinesIfNotFound(fs, "/f.txt", "1", "2"))
	a.Equal("1\n2", ufs.ReadText(fs, "/f.txt"))
}

func Test_TempFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	actual := ufs.TempFile(fs, "xyz")
	a.True(strings.Contains(actual, "xyz"))
	a.NotEqual("xyz", actual)
}

func Test_ExpandHomePath_happy(t *testing.T) {
	a := require.New(t)
	a.Equal("none", ufs.ExpandHomePath("none"))

	a.Equal(ufs.UserHomeDir(), ufs.ExpandHomePath("~"))
}
