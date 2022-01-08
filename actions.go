package main

import (
	"compress/gzip"
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

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// determine the directory of the file to be archived in
	// relation to the root. filepath.* is cross-platform
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	// construct the path to the archived file's location
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// open the target file with r/w perms
	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	// open the source file
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	// create new zip writer
	zw := gzip.NewWriter(out)

	// store the source file name in the compressed file
	zw.Name = filepath.Base(path)

	// copy the data into the compressed format
	if _, err := io.Copy(zw, in); err != nil {
		return err
	}

	// don't defer .Close() bc we want to return the error
	if err := zw.Close(); err != nil {
		return err
	}

	return out.Close()
}
