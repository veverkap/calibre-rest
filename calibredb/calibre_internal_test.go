package calibredb

import (
	"testing"
)

func TestFiltered(t *testing.T) {
	tests := []struct {
		name    string
		output  []byte
		isError bool
		want    string
	}{
		{
			name:    "Empty output",
			output:  []byte(""),
			isError: false,
			want:    "",
		},
		{
			name:    "Simple output without error",
			output:  []byte("Line 1\nLine 2\nLine 3"),
			isError: false,
			want:    "Line 1\nLine 2\nLine 3",
		},
		{
			name:    "Output with blank lines",
			output:  []byte("Line 1\n\nLine 2\n\nLine 3"),
			isError: false,
			want:    "Line 1\nLine 2\nLine 3",
		},
		{
			name:    "Output with Integration status line",
			output:  []byte("Line 1\nIntegration status: False\nLine 2"),
			isError: false,
			want:    "Line 1\nLine 2",
		},
		{
			name:    "Error output returns last line only",
			output:  []byte("Traceback line 1\nTraceback line 2\napsw.ConstraintError: UNIQUE constraint failed"),
			isError: true,
			want:    "apsw.ConstraintError: UNIQUE constraint failed",
		},
		{
			name:    "Error output with blank lines and Integration status",
			output:  []byte("Traceback line 1\n\nTraceback line 2\nIntegration status: False\napsw.ConstraintError: UNIQUE constraint failed"),
			isError: true,
			want:    "apsw.ConstraintError: UNIQUE constraint failed",
		},
		{
			name:    "Complex error with multiple Integration status lines",
			output:  []byte("Traceback line 1\nIntegration status: True\nTraceback line 2\nIntegration status: False\nError message"),
			isError: true,
			want:    "Error message",
		},
		{
			name:    "Error output with only blank lines",
			output:  []byte("\n\n\n"),
			isError: true,
			want:    "",
		},
		{
			name: "Realistic error traceback",
			output: []byte(`Traceback (most recent call last):
  File "runpy.py", line 198, in _run_module_as_main
  File "runpy.py", line 88, in _run_code
  File "calibre/db/cli/cmd_add_custom_column.py", line 81, in main
apsw.ConstraintError: UNIQUE constraint failed: custom_columns.label
Integration status: False`),
			isError: true,
			want:    "apsw.ConstraintError: UNIQUE constraint failed: custom_columns.label",
		},
		{
			name: "Non-error output with traceback-like content",
			output: []byte(`Traceback (most recent call last):
  File "runpy.py", line 198, in _run_module_as_main
  File "runpy.py", line 88, in _run_code
apsw.ConstraintError: UNIQUE constraint failed: custom_columns.label
Integration status: False`),
			isError: false,
			want:    "Traceback (most recent call last):\n  File \"runpy.py\", line 198, in _run_module_as_main\n  File \"runpy.py\", line 88, in _run_code\napsw.ConstraintError: UNIQUE constraint failed: custom_columns.label",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filtered(tt.output, tt.isError)
			if got != tt.want {
				t.Errorf("filtered() = %q, want %q", got, tt.want)
			}
		})
	}
}
