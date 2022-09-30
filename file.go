package ufs

import (
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/fastgh/go-comm/v2"
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
	DownloadP() Content
	Download() (Content, error)
}

type (
	CredentialsT = models.Credentials
	Credentials  = *CredentialsT
)

type (
	ContentT = models.RemoteFileContent
	Content  = *ContentT
)

func NewFileP(afs afero.Fs, url string, credentials Credentials, timeout time.Duration) File {
	r, err := NewFile(afs, url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func NewFile(afs afero.Fs, url string, credentials Credentials, timeout time.Duration) (File, error) {
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

func DownloadBytesP(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) []byte {
	r, err := DownloadBytes(fs, url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func DownloadBytes(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) ([]byte, error) {
	f, err := NewFile(fs, url, credentials, timeout)
	if err != nil {
		return nil, err
	}

	c, err := f.Download()
	if err != nil {
		return nil, err
	}

	blob := c.Blob
	defer blob.Close()

	return comm.ReadBytes(blob)
}

func DownloadTextP(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) string {
	r, err := DownloadText(fs, url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func DownloadText(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) (string, error) {
	bytes, err := DownloadBytes(fs, url, credentials, timeout)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MapFromYamlFileP(fs afero.Fs, path string, envsubt bool) map[string]any {
	r, err := MapFromYamlFile(fs, path, envsubt)
	if err != nil {
		panic(err)
	}
	return r
}

func MapFromYamlFile(fs afero.Fs, path string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromYamlFile(fs, path, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func FromYamlFileP(fs afero.Fs, path string, envsubt bool, result any) {
	if err := FromYamlFile(fs, path, envsubt, result); err != nil {
		panic(err)
	}
}

func FromYamlFile(fs afero.Fs, path string, envsubt bool, result any) error {
	yamlText, err := ReadText(fs, path)
	if err != nil {
		return err
	}

	if envsubt {
		yamlText, err = comm.EnvSubst(yamlText, nil)
		if err != nil {
			return err
		}

	}

	if err := yaml.Unmarshal([]byte(yamlText), result); err != nil {
		return errors.Wrapf(err, "failed to parse yaml file: %s\n\n%s", path, yamlText)
	}
	return nil
}
