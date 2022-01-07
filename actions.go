package main

import (
	"fmt"
	"io"
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

func delFile(path string) error {
	// os.Remove() returns an error that will bubble up and return
	return os.Remove(path)
}
