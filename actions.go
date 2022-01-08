package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// filterOut checks if the path must be filtered out from the results
// based on whether the path points to a dir, the file size is less than
// the min file size provided by the user, or the file extension doesn't
// match the extension provided by the user
func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}
	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

// Returns an error if the path cannot be removed, logs the deleted
// path and returns nil if it is deleted.
func delFile(path string, delLogger *log.Logger) error {
	// os.Remove() returns an error that will bubble up and return
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}
