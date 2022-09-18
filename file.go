package ufs

import (
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/fastgh/go-comm"
	"github.com/goodsru/go-universal-network-adapter/models"
	"github.com/goodsru/go-universal-network-adapter/services"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	HTTP  = services.HTTP + "://"
	HTTPS = services.HTTPS + "://"
	FTP   = services.FTP + "://"
	FTPS  = services.FTPS + "://"
	SFTP  = services.SFTP + "://"
	S3    = services.S3 + "://"
	FILE  = "file://"
)

type File interface {
	Name() string
	Dir() string
	Url() string
	Protocol() string
	URL() *url.URL
	Credentials() Credentials
	Timeout() time.Duration
	Download() Content
}

type CredentialsT = models.Credentials
type Credentials = *CredentialsT

type ContentT = models.RemoteFileContent
type Content = *ContentT

func NewFile(afs afero.Fs, url string, credentials Credentials, timeout time.Duration) File {
	if IsRemote(url) {
		return NewRemoteFile(url, credentials, timeout)
	}
	return NewAferoFile(afs, url, credentials, timeout)
}

func IsFileProtocol(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), FILE)
}

func IsRemote(url string) bool {
	lc := strings.ToLower(url)

	if strings.HasPrefix(lc, HTTP) ||
		strings.HasPrefix(lc, HTTPS) ||
		strings.HasPrefix(lc, FTP) ||
		strings.HasPrefix(lc, FTPS) ||
		strings.HasPrefix(lc, SFTP) ||
		strings.HasPrefix(lc, S3) {
		return true
	}
	return false
}

func WorkDir(url string, defaultDir string) string {
	if IsRemote(url) {
		return defaultDir
	}
	if IsFileProtocol(url) {
		url = url[len(FILE):]
	}

	r := filepath.Dir(url)
	if r == "." {
		return defaultDir
	}
	if !filepath.IsAbs(url) {
		r = filepath.Join(defaultDir, r)
	}
	return r
}

func ShortDescription(url string) string {
	r := url

	protocol := ""
	posOfProtocolSep := strings.Index(strings.ToLower(r), "://")
	if posOfProtocolSep >= 0 {
		protocol = r[:posOfProtocolSep+3]
		r = r[posOfProtocolSep+len("://"):]
	}

	lem := len(r)
	if lem <= 3+8+1+5 /* ...12345678.hosts */ {
		return protocol + r
	}

	return protocol + r[:3] + "..." + r[lem-(8+1+5): /* 12345678.hosts */]
}

func DownloadBytes(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) []byte {
	f := NewFile(fs, url, credentials, timeout)
	blob := f.Download().Blob
	defer blob.Close()

	return comm.ReadBytes(blob)
}

func DownloadText(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) string {
	bytes := DownloadBytes(fs, url, credentials, timeout)
	return string(bytes)
}

func MapFromYamlFile(fs afero.Fs, path string) map[string]any {
	r := map[string]any{}
	FromYamlFile(fs, path, &r)
	return r
}

func FromYamlFile(fs afero.Fs, path string, result any) {
	yamlText := ReadText(fs, path)
	yamlText = comm.EnvSubst(yamlText, nil)

	if err := yaml.Unmarshal([]byte(yamlText), result); err != nil {
		panic(errors.Wrapf(err, "failed to parse yaml file: %s\n\n%s", path, yamlText))
	}
}
