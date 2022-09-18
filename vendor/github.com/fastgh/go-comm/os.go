package comm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/a8m/envsubst/parse"
	"github.com/joho/godotenv"
	plog "github.com/phuslu/log"
	"github.com/pkg/errors"
)

// OSType ...
type OSType int

const (
	// AllOSType ...
	AllOSType OSType = iota

	// Windows ...
	Windows

	// Darwin ...
	Darwin

	// Linux ...
	Linux
)

// ParseOSType ...
func ParseOSType(s string) OSType {
	if s == "all" {
		return AllOSType
	}
	if s == "windows" {
		return Windows
	}
	if s == "darwin" {
		return Darwin
	}
	if s == "linux" {
		return Linux
	}
	panic(fmt.Errorf("unknown OS type: '%s'", s))
}

// BuildOSType ...
func BuildOSType(i int) OSType {
	if int(AllOSType) == i {
		return AllOSType
	}
	if int(Windows) == i {
		return Windows
	}
	if int(Darwin) == i {
		return Darwin
	}
	if int(Linux) == i {
		return Linux
	}
	panic(fmt.Errorf("unknown OS type: '%v'", i))
}

// DefaultOSType ...
func DefaultOSType() OSType {
	return ParseOSType(runtime.GOOS)
}

func IsWindows() bool {
	return DefaultOSType() == Windows
}

func IsDarwin() bool {
	return DefaultOSType() == Darwin
}

func EnvironMap(overrides map[string]string) map[string]string {
	envs := JoinedLines(os.Environ()...)
	r, err := godotenv.Unmarshal(envs)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse OS environments"))
	}

	if len(overrides) > 0 {
		for k, v := range overrides {
			r[k] = v
		}
	}
	return r
}

func EnvironList(overrides map[string]string) []string {
	envs := EnvironMap(overrides)
	r := make([]string, 0, len(envs)+len(overrides))
	for k, v := range envs {
		r = append(r, k+"="+v)
	}
	return r
}

func EnvSubst(input string, env map[string]string) string {
	restr := parse.Restrictions{NoUnset: false, NoEmpty: false}
	parser := parse.New("tmp", EnvironList(EnvironMap(env)), &restr)

	var r string
	var err error
	if r, err = parser.Parse(input); err != nil {
		panic(errors.Wrapf(err, "failed to envsubst the text: %s", input))
	}
	return r
}

func EnvSubstSlice(inputs []string, env map[string]string) []string {
	r := make([]string, 0, len(inputs))
	for _, s := range inputs {
		r = append(r, EnvSubst(s, env))
	}
	return r
}

func IsTerminal() bool {
	return plog.IsTerminal(os.Stdout.Fd())
}

func Executable() string {
	r, err := os.Executable()
	if err != nil {
		panic(errors.Wrap(err, "failed to get the path name of the executable file"))
	}
	r, err = filepath.EvalSymlinks(r)
	if err != nil {
		panic(errors.Wrapf(err, "failed to evaluate the symbol linke of the executable file: %s", r))
	}
	return r
}

func WorkingDirectory() string {
	r, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "failed to get working directory"))
	}
	return r
}
