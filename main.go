package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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
	// log destination writer
	// accepts a file in main() or a buffer for testing
	wLog io.Writer
	// archive directory
	archive string
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")
	// if no log file is provided, use STDOUT
	logFile := flag.String("log", "", "Log deletes to this file")
	// Action options
	list := flag.Bool("list", false, "List files only")
	archive := flag.String("archive", "", "Archive directory")
	del := flag.Bool("del", false, "Delete files")
	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	// write to stdout by default
	var (
		f   = os.Stdout
		err error
	)

	// if a value was provided for the logFile, open it and store in f
	// no need to test this because it uses std lib funcs that are tested
	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext:     *ext,
		size:    *size,
		list:    *list,
		del:     *del,
		wLog:    f,
		archive: *archive,
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
	// new logger using wLog, prefix, and std log flags like date and time
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)
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

			// Archive files and continue if successful
			if cfg.archive != "" {
				// return if there is an error only
				if err := archiveFile(cfg.archive, root, path); err != nil {
					return err
				}
			}

			// Delete files
			if cfg.del {
				return delFile(path, delLogger)
			}

			// List is the default option if nothing else was set
			return listFile(path, out)
		})
}
