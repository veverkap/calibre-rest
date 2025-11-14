package calibredb

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/samber/lo"
)

type Calibre struct {
	CalibreDBLocation string `json:"calibredb_location,omitempty"`
	LibraryPath       string `json:"library-path,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	Timeout           string `json:"timeout,omitempty"`
	OnError           func(error)

	validate *validator.Validate
}

type CalibreOption func(*Calibre)

func WithCalibreDBLocation(path string) CalibreOption {
	return func(c *Calibre) {
		c.CalibreDBLocation = path
	}
}

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
	c.validate = validator.New(validator.WithRequiredStructEnabled())
	return c
}

func (c *Calibre) Version() string {
	if out, err := c.run("--version"); err != nil {
		return err.Error()
	} else {
		return out
	}
}

func (c *Calibre) Help() string {
	if out, err := c.run("--help"); err != nil {
		return err.Error()
	} else {
		return out
	}
}

func (c *Calibre) run(argv ...string) (string, error) {
	argv = append(argv, "--with-library="+c.LibraryPath)
	out, err := exec.Command(c.CalibreDBLocation, argv...).CombinedOutput()
	if err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		if out != nil {
			// this is a stacktrace followed by the actual error message. We want to extract only the actual error message.
			return "", errors.New(filtered(out, true))
		}
		return "", errors.New(err.Error())
	}
	return filtered(out, false), nil
}

func filtered(output []byte, isError bool) string {
	// The format of the error is a traceback followed by the actual error message. We want to extract only the actual error message.
	// Example:
	// 	Traceback (most recent call last):
	//   File "runpy.py", line 198, in _run_module_as_main
	//   File "runpy.py", line 88, in _run_code
	//   File "site.py", line 42, in <module>
	//   File "site.py", line 38, in main
	//   File "calibre/db/cli/main.py", line 253, in main
	//   File "calibre/db/cli/main.py", line 40, in run_cmd
	//   File "calibre/db/cli/cmd_add_custom_column.py", line 81, in main
	//   File "calibre/db/cli/cmd_add_custom_column.py", line 72, in do_add_custom_column
	//   File "calibre/db/legacy.py", line 812, in create_custom_column
	//   File "calibre/db/cache.py", line 86, in call_func_with_lock
	//   File "calibre/db/cache.py", line 2669, in create_custom_column
	//   File "calibre/db/backend.py", line 1244, in create_custom_column
	//   File "calibre/db/backend.py", line 1171, in execute
	//   File "src/cursor.c", line 189, in resetcursor
	// apsw.ConstraintError: UNIQUE constraint failed: custom_columns.label
	// Integration status: False
	outputstring := string(output)
	outputLines := strings.Split(outputstring, "\n")
	// remove any lines that start with "Integration status:" or that are blank
	filteredLines := make([]string, 0, len(outputLines))
	for _, line := range outputLines {
		if line != "" && !strings.HasPrefix(line, "Integration status:") {
			filteredLines = append(filteredLines, line)
		}
	}
	if isError {
		// return the last line only
		if len(filteredLines) > 0 {
			return filteredLines[len(filteredLines)-1]
		}
	}
	return strings.Join(filteredLines, "\n")
}
