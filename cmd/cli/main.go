package main

import (
	"fmt"
	"os"

	"github.com/veverkap/calibre-rest/calibredb"
)

func main() {
	// get temp directory
	tempDir := os.TempDir()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)
	out := c.Help()
	fmt.Println(out)

	output, err := c.AddCustomColumn(
		calibredb.AddCustomColumnOptions{
			Label:    "sdfassaddafsdsdfdsfsd",
			Name:     "My Column",
			Datatype: "text",
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}

// c.Help()
// c.SetMetadata(
// 	calibredb.SetMetadataOptions{
// 		BookId:     "1",
// 		Path:       "fsad",
// 		ListFields: lo.ToPtr(false),
// 	},
// )
// output, err := c.AddCustomColumn(
// 	calibredb.AddCustomColumnOptions{
// 		Label:    "sdfassaddafsdsdfdsfsd",
// 		Name:     "My Column",
// 		Datatype: "text",
// 	},
// )
// fmt.Println(output, err)
// os.Exit(1)

// func test() {
// 	c := calibredb.NewCalibre(calibredb.WithLibraryPath("/Users/veverkap/Code/personal/calibre-rest"))

// 	s := c.List(
// 		calibredb.ListOptions{
// 			Ascending:  true,
// 			Fields:     "author_sort, authors, comments",
// 			ForMachine: lo.ToPtr(false),
// 			Limit:      2,
// 		},
// 	)
// 	fmt.Println(s)
// }
