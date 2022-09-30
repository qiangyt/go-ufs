package test

import (
	"io"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/fastgh/go-ufs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AferoFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	ufs.WriteTextP(fs, "/hello/world/test.txt", "hi")

	cred := &ufs.CredentialsT{User: "u"}
	a1 := ufs.NewAferoFileP(fs, "/hello/world/test.txt", cred, 3*time.Second)
	a2 := ufs.NewAferoFileP(fs, "file:///hello/world/test.txt", cred, 3*time.Second)

	a.Same(fs, a1.Fs())
	a.Same(fs, a2.Fs())
	a.Equal("test.txt", a1.Name())
	a.Equal("test.txt", a2.Name())
	a.Equal(filepath.Join("/", "hello", "world"), a1.Dir())
	a.Equal(filepath.Join("/", "hello", "world"), a2.Dir())
	a.Equal("file", a1.Protocol())
	a.Equal("file", a2.Protocol())
	a.Equal("file:///hello/world/test.txt", a1.Url())
	a.Equal("file:///hello/world/test.txt", a2.Url())
	a.Same(cred, a1.Credentials())
	a.Same(cred, a2.Credentials())
	a.Equal(3*time.Second, a1.Timeout())
	a.Equal(3*time.Second, a2.Timeout())

	tUrl, _ := url.Parse("file:///hello/world/test.txt")
	a.Equal(tUrl, a1.URL())
	a.Equal(tUrl, a2.URL())

	c1 := a1.DownloadP()
	a.Equal("test.txt", c1.Name)
	a.Equal("/hello/world/test.txt", c1.Path)

	txt1, _ := io.ReadAll(c1.Blob)
	defer c1.Blob.Close()
	a.Equal("hi", string(txt1))

	c2 := a2.DownloadP()
	a.Equal("test.txt", c2.Name)
	a.Equal("/hello/world/test.txt", c2.Path)

	txt2, _ := io.ReadAll(c2.Blob)
	defer c2.Blob.Close()
	a.Equal("hi", string(txt2))
}
