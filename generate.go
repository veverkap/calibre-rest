//go:build ignore

// To run: go run generate_calibredb_wrappers.go
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/samber/lo"
)

var cmds = []string{
	"add_custom_column",
	"add_format",
	"add",
	"backup_metadata",
	"catalog",
	"check_library",
	"clone",
	"custom_columns",
	"embed_metadata",
	"export",
	"fts_index",
	"fts_search",
	"list_categories",
	"list",
	"remove_custom_column",
	"remove_format",
	"remove",
	"restore_database",
	"saved_searches",
	"search",
	"set_custom",
	"set_metadata",
	"show_metadata",
}

func main() {
	// c := calibredb.NewCalibre(calibredb.WithLibraryPath("/Users/veverkap/Code/personal/calibre-rest"))
	// help := c.ListHelp()
	// fmt.Println("Calibre Help:", help)
	for _, cmd := range cmds {
		var out bytes.Buffer
		out.WriteString(`package calibredb
		
	`)
		fmt.Println("Processing cmd", cmd)
		flags, description := parseCommandHelp(cmd)
		// write description to a comment in the file
		pascalCmd := lo.PascalCase(cmd)
		structName := pascalCmd + "Options"
		out.WriteString("// " + strings.ReplaceAll((description), "\n", "\n// ") + "\n")
		out.WriteString("type " + structName + " struct {\n")
		for _, flag := range flags {
			fieldName := lo.PascalCase(flag)
			out.WriteString(fmt.Sprintf("\t%s string `json:\"%s,omitempty\"`\n", fieldName, flag))
		}
		out.WriteString("}\n")
		out.WriteString("\n")
		out.WriteString("func (c *Calibre) " + pascalCmd + "Help() string {\n")
		out.WriteString(fmt.Sprintf("\treturn c.run(\"%s\", \"-h\")\n", cmd))
		out.WriteString("}\n")
		out.WriteString("\n")
		out.WriteString("func (c *Calibre) " + pascalCmd + "(opts " + structName + ", args ...string) string {\n")
		out.WriteString(fmt.Sprintf("\treturn \"%s\"\n", cmd))
		out.WriteString("}\n")
		err := os.WriteFile(fmt.Sprintf("calibredb/%s.go", cmd), out.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
		out.Reset()
	}
}

// GLOBAL OPTIONS:
//
//	--library-path=LIBRARY_PATH, --with-library=LIBRARY_PATH
//	-h, --help          show this help message and exit
//	--version           show program's version number and exit
//	--username=USERNAME
//	--password=PASSWORD
//	--timeout=TIMEOUT   The timeout, in seconds, when connecting to a calibre
var skippedGlobalOptions = map[string]struct{}{
	"library-path": {},
	"with-library": {},
	"help":         {},
	"version":      {},
	"username":     {},
	"password":     {},
	"timeout":      {},
}

func parseCommandHelp(cmd string) ([]string, string) {
	result := make(map[string]bool)
	out, err := exec.Command("/Applications/calibre.app/Contents/MacOS/calibredb", cmd, "-h").CombinedOutput()

	if err != nil {
		// Skip commands that don't support -h
		return nil, ""
	}

	// os.WriteFile("docs/"+cmd+".md", out, 0666)
	stringout := string(out)
	// skip the first line if it contains "Integration status: "
	if strings.HasPrefix(stringout, "Integration status: ") {
		stringout = strings.Join(strings.SplitN(stringout, "\n", 2)[1:], "\n")
	}
	lines := strings.Split(stringout, "\n")
	flagRe := regexp.MustCompile(`\s+--([a-zA-Z0-9_-]+)([=\s][^ ]*)?\s+(.*)$`)
	for _, l := range lines {
		if m := flagRe.FindStringSubmatch(l); len(m) >= 4 {
			flagName := strings.TrimSpace(m[1])
			if _, ok := skippedGlobalOptions[flagName]; ok {
				continue
			}
			result[flagName] = true
		}
	}
	sortedResult := lo.Keys(result)
	sort.Strings(sortedResult)
	return sortedResult, stringout
}
