package test

import (
	"io"
	"testing"
	"time"

	"github.com/fastgh/go-ufs"
	"github.com/stretchr/testify/require"
)

func Test_RemoteFile_happy(t *testing.T) {
	a := require.New(t)

	// cred := &CredentialsT{}
	actual := ufs.NewRemoteFile("https://mirror.sjtu.edu.cn/debian/README.mirrors.txt", nil, 10*time.Second)

	a.Equal("README.mirrors.txt", actual.Name())
	a.Equal("/debian", actual.Dir())
	a.Equal("https", actual.Protocol())
	a.Equal("https://mirror.sjtu.edu.cn/debian/README.mirrors.txt", actual.Url())
	// a.Equal(cred, actual.Credentials())
	a.Equal(10*time.Second, actual.Timeout())

	c := actual.Download()
	a.Equal("README.mirrors.txt", c.Name)
	a.Contains(c.Path, "README.mirrors.txt")

	txt, _ := io.ReadAll(c.Blob)
	defer c.Blob.Close()

	a.Equal("The list of Debian mirror sites is available here: https://www.debian.org/mirror/list\n", string(txt))
}
