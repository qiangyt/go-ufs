package test

import (
	"path/filepath"
	"testing"

	"github.com/fastgh/go-comm/v2"
	"github.com/fastgh/go-ufs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_WorkDir_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("defaultDir", ufs.WorkDir("http://test", "defaultDir"))
	a.Equal("defaultDir", ufs.WorkDir("file://hello.txt", "defaultDir"))
	a.Equal("defaultDir", ufs.WorkDir("hello.txt", "defaultDir"))
	a.Equal(filepath.Join("defaultDir", "home"), ufs.WorkDir("home/hello.txt", "defaultDir"))

	if comm.IsWindows() {
		a.Equal("c:\\home", ufs.WorkDir("file://c:/home/hello.txt", "defaultDir"))
		a.Equal("c:\\home", ufs.WorkDir("c:/home/hello.txt", "defaultDir"))

		a.Equal("C:\\", ufs.WorkDir("C:\\hello.txt", "defaultDir"))
	} else {
		a.Equal("/home", ufs.WorkDir("file:///home/hello.txt", "defaultDir"))
		a.Equal("/home", ufs.WorkDir("/home/hello.txt", "defaultDir"))
	}
}

func Test_IsRemote_happy(t *testing.T) {
	a := require.New(t)

	a.True(ufs.IsRemote("http://test.local"))
	a.True(ufs.IsRemote("HTTP://test.local"))
	a.True(ufs.IsRemote("https://test.local"))
	a.True(ufs.IsRemote("HTTPS://test.local"))
	a.True(ufs.IsRemote("HTTPs://test.local"))

	a.True(ufs.IsRemote("ftp://test.local"))
	a.True(ufs.IsRemote("FTP://test.local"))
	a.True(ufs.IsRemote("ftps://test.local"))
	a.True(ufs.IsRemote("FTPS://test.local"))
	a.True(ufs.IsRemote("FTPs://test.local"))

	a.True(ufs.IsRemote("sftp://test.local"))
	a.True(ufs.IsRemote("SFTP://test.local"))

	a.True(ufs.IsRemote("s3://test.local"))
	a.True(ufs.IsRemote("S3://test.local"))
}

func Test_NewFile(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	fRemote := ufs.NewFileP(nil, "https://google.com", nil, 0)
	_, isRemoteFile := fRemote.(ufs.RemoteFile)
	a.True(isRemoteFile)

	fLocal := ufs.NewFileP(afs, "file://test.txt", nil, 0)
	_, isLocalFile := fLocal.(ufs.AferoFile)
	a.True(isLocalFile)

	fLocal = ufs.NewFileP(afs, "test.txt", nil, 0)
	_, isLocalFile = fLocal.(ufs.AferoFile)
	a.True(isLocalFile)
}

func Test_ShortDescription_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("", ufs.ShortDescription(""))

	a.Equal("AB/12345678.hosts", ufs.ShortDescription("AB/12345678.hosts"))
	a.Equal("ABC...12345678.hosts", ufs.ShortDescription("ABC/12345678.hosts"))

	a.Equal("https://AB/12345678.hosts", ufs.ShortDescription("https://AB/12345678.hosts"))
	a.Equal("SFTP://ABC...12345678.hosts", ufs.ShortDescription("SFTP://ABC/12345678.hosts"))
}

func Test_DownloadText_happy(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	ufs.WriteTextP(afs, "test.txt", "Test_DownloadText_happy")

	actual := ufs.DownloadTextP(afs, "test.txt", nil, 0)
	a.Equal("Test_DownloadText_happy", actual)
}

func Test_DownloadText_Remote(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	actual := ufs.DownloadTextP(afs, "https://mirror.sjtu.edu.cn/debian/README.mirrors.txt", nil, 0)
	a.Equal("The list of Debian mirror sites is available here: https://www.debian.org/mirror/list\n", actual)
}

func Test_YamlFileToMap_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.WriteTextP(fs, "test.yaml", `k: v`)

	configMap := ufs.MapFromYamlFileP(fs, "test.yaml", false)

	a.Len(configMap, 1)
	a.Equal("v", configMap["k"])
}
