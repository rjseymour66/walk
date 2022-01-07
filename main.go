package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// config contains values that would otherwise be passed to
// run as params. This is to clean up the function definition
type config struct {
	// extension to filter out
	ext string
	// min file size
	size int64
	// list files
	list bool
	// delete files
	del bool
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")
	// Action options
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run executes the program
// root: dir to start the search
// out: output destintation. Allows you to print results to STDOUT or to a
//      bytes.Buffer when testing.
func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {

			// Can func walk this directory
			if err != nil {
				return err
			}

			// Should the current directory or file be filtered out. If so,
			// returns nil to skip the rest of the function and processes the
			// next file or dir
			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}

			// If list was explicitly set, don't do anything else
			if cfg.list {
				return listFile(path, out)
			}

			// Delete files
			if cfg.del {
				return delFile(path)
			}

			// List is the default option if nothing else was set
			return listFile(path, out)
		})
}
