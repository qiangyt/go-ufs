package comm

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func DefaultOutput() io.Writer {
	if IsTerminal() {
		return os.Stdout
	} else {
		return io.Discard
	}
}

// ReadBytes ...
func ReadBytes(reader io.Reader) []byte {
	r, err := io.ReadAll(reader)
	if err != nil {
		panic(errors.Wrap(err, "failed to read from Reader"))
	}
	return r
}

// ReadText ...
func ReadText(reader io.Reader) string {
	return string(ReadBytes(reader))
}

func ReadLines(reader io.Reader) []string {
	r := make([]string, 0, 32)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		r = append(r, line)
	}

	return r
}

func Text2Lines(text string) []string {
	rdr := strings.NewReader(text)
	return ReadLines(rdr)
}
