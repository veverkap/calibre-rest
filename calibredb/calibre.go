package calibredb

import (
	"os/exec"

	_ "github.com/samber/lo"
)

type Calibre struct {
	LibraryPath string `json:"library-path,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Timeout     string `json:"timeout,omitempty"`
	OnError     func(error)
}

type CalibreOption func(*Calibre)

func WithLibraryPath(path string) CalibreOption {
	return func(c *Calibre) {
		c.LibraryPath = path
	}
}

func WithUsername(username string) CalibreOption {
	return func(c *Calibre) {
		c.Username = username
	}
}

func WithPassword(password string) CalibreOption {
	return func(c *Calibre) {
		c.Password = password
	}
}

func WithTimeout(timeout string) CalibreOption {
	return func(c *Calibre) {
		c.Timeout = timeout
	}
}

func WithOnError(onError func(error)) CalibreOption {
	return func(c *Calibre) {
		c.OnError = onError
	}
}

func NewCalibre(opts ...CalibreOption) *Calibre {
	c := &Calibre{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Calibre) Version() string {
	return c.run("--version")
}

func (c *Calibre) Help() string {
	return c.run("--help")
}

func (c *Calibre) run(argv ...string) string {
	argv = append(argv, "--with-library="+c.LibraryPath)
	out, err := exec.Command("/Applications/calibre.app/Contents/MacOS/calibredb", argv...).CombinedOutput()
	if err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		if out != nil {
			return string(out)
		}
		return err.Error()
	}
	return string(out)
}

// cmd := exec.Command("/Applications/calibre.app/Contents/MacOS/calibredb", "list", "--limit=5", "-f", "title,authors,author_sort,tags,isbn", "--for-machine", "--with-library=.")
