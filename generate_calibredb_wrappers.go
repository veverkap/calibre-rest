//go:build ignore

// To run: go run generate_calibredb_wrappers.go
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Flag struct {
	Name        string
	Description string
}

type Command struct {
	Name        string
	Description string
	Flags       []Flag
}

var cmds = []Command{
	{
		Name:        "list",
		Description: "List all books in the calibre database.",
	},
	{
		Name:        "add",
		Description: "Add a book to the calibre database.",
	},
	{
		Name:        "remove",
		Description: "Remove a book from the calibre database.",
	},
	{
		Name:        "show_metadata",
		Description: "",
	},
	{
		Name:        "set_metadata",
		Description: "",
	},
	{
		Name:        "export",
		Description: "",
	},
	{
		Name:        "catalog",
		Description: "",
	},
	{
		Name:        "saved_searches",
		Description: "",
	},
	{
		Name:        "add_custom_column",
		Description: "",
	},
	{
		Name:        "custom_columns",
		Description: "",
	},
	{
		Name:        "remove_custom_column",
		Description: "",
	},
	{
		Name:        "set_custom",
		Description: "",
	},
	{
		Name:        "restore_database",
		Description: "",
	},
	{
		Name:        "check_library",
		Description: "",
	},
	{
		Name:        "list_categories",
		Description: "",
	},
	{
		Name:        "backup_metadata",
		Description: "",
	},
	{
		Name:        "clone",
		Description: "",
	},
	{
		Name:        "embed_metadata",
		Description: "",
	},
	{
		Name:        "search",
		Description: "",
	},
	{
		Name:        "fts_index",
		Description: "",
	},
	{
		Name:        "fts_search",
		Description: "",
	},
}

func main() {
	completeCmds := cmds
	for i, cmd := range cmds {
		completeCmds[i].Flags = parseCommandHelp(cmd.Name)
	}

	// Step 3: Generate Go code
	var out bytes.Buffer

	for _, cmd := range completeCmds {
		out.WriteString(`package calibrewrap
	`)

		fmt.Println("Processing cmd:", cmd.Name)
		funcName := cases.Upper(language.English).String(cmd.Name)
		out.WriteString(fmt.Sprintf("\n// Calibredb%s executes `calibredb %s`.\n", funcName, cmd.Name))
		out.WriteString("//\n")
		if cmd.Description != "" {
			out.WriteString("// " + cmd.Description + "\n")
		}
		hasFlags := len(cmd.Flags) > 0
		if hasFlags {
			out.WriteString("//\n// Flags:\n")
			for _, f := range cmd.Flags {
				out.WriteString(fmt.Sprintf("//   --%s: %s\n", f.Name, f.Description))
			}

			out.WriteString(fmt.Sprintf("\ntype Calibredb%sFlags struct {\n", funcName))
			for _, f := range cmd.Flags {
				fieldName := cases.Title(language.English).String(strings.ReplaceAll(f.Name, "-", " "))
				fieldName = strings.ReplaceAll(fieldName, " ", "")
				out.WriteString(fmt.Sprintf("\t%s string // %s\n", fieldName, f.Description))
			}
			out.WriteString("}\n")
		}

		// // Build the function body
		// if hasFlags {
		// 	out.WriteString(fmt.Sprintf("\nfunc Calibredb%s(flags Calibredb%sFlags, args ...string) Result {\n", funcName, funcName))
		// 	out.WriteString("\tvar argv []string\n")
		// 	out.WriteString(fmt.Sprintf("\targv = append(argv, \"%s\")\n", cmd.Name))
		// 	out.WriteString("\t// Process flags\n")
		// 	for _, f := range cmd.Flags {
		// 		fieldName := cases.Title(language.English).String(strings.ReplaceAll(f.Name, "-", " "))
		// 		fieldName = strings.ReplaceAll(fieldName, " ", "")
		// 		out.WriteString(fmt.Sprintf("\tif flags.%s != \"\" {\n", fieldName))
		// 		out.WriteString(fmt.Sprintf("\t\targv = append(argv, fmt.Sprintf(\"--%s=%s\", \"%s\", flags.%s))\n", f.Name, "%s", f.Name, fieldName))
		// 		out.WriteString("\t}\n")
		// 	}
		// 	out.WriteString("\targv = append(argv, args...)\n")
		// 	out.WriteString("\treturn run(argv...)\n}\n")
		// 	continue
		// }

		// // No flags
		// out.WriteString(fmt.Sprintf("func Calibredb%s(flags map[string]string, args ...string) Result {\n", funcName))
		// out.WriteString("\tvar argv []string\n")
		// out.WriteString(fmt.Sprintf("\targv = append(argv, \"%s\")\n", cmd.Name))
		// out.WriteString("\tfor k, v := range flags {\n")
		// out.WriteString("\t\tif v == \"\" {\n")
		// out.WriteString("\t\t\targv = append(argv, fmt.Sprintf(\"--%s\", k))\n")
		// out.WriteString("\t\t} else {\n")
		// out.WriteString("\t\t\targv = append(argv, fmt.Sprintf(\"--%s=%s\", k, v))\n")
		// out.WriteString("\t\t}\n")
		// out.WriteString("\t}\n")
		// out.WriteString("\targv = append(argv, args...)\n")
		// out.WriteString("\treturn run(argv...)\n}\n")
		err := os.WriteFile(fmt.Sprintf("calibrewrap/calibredb_%s.go", cmd.Name), out.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
		out.Reset()
	}

	err := os.WriteFile("calibrewrap/calibredb_gen.go", out.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ Generated calibredb_gen.go with %d commands\n", len(cmds))
}

func parseCommandHelp(cmd string) []Flag {
	out, err := exec.Command("/Applications/calibre.app/Contents/MacOS/calibredb", cmd, "-h").CombinedOutput()
	if err != nil {
		// Skip commands that don't support -h
		return nil
	}
	os.WriteFile("docs/"+cmd+".md", out, 0666)
	lines := strings.Split(string(out), "\n")
	flagRe := regexp.MustCompile(`\s+--([a-zA-Z0-9_-]+)([=\s][^ ]*)?\s+(.*)$`)
	var flags []Flag
	var uniqueFlags = make(map[string]bool)
	for _, l := range lines {
		if m := flagRe.FindStringSubmatch(l); len(m) >= 4 {
			if ok, _ := uniqueFlags[m[1]]; ok {
				continue
			}
			uniqueFlags[m[1]] = true
			flags = append(flags, Flag{
				Name:        strings.TrimSpace(m[1]),
				Description: strings.TrimSpace(m[3]),
			})
		}
	}
	return flags
}
