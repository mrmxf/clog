//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

// Package care manages the embedded files of clog
// if you are forking clog, then you probably want to override this package
// to include ONLY the files that you want in your own fork.

package core

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

//variable CoreFs is the exported file system

//go:embed all:tpl
//go:embed core.clog.yaml
var CoreFs embed.FS

// var discardPrefix is for discarding prefixes that are added when using
// golang embed and the folder being included should not be visible to the end
// user. CAUTION - this will be problematic if the prefix is also the same
// as the name of a top level folder - don't do that.

var discardPrefix = "core/"

// Clean will remove any prefix from the file path and do other path cleaning
// functions to improve extensibility of the embed FS
func Clean(dirtyPath string) (cleanPath string) {
	cleanPath = dirtyPath
	if strings.HasPrefix(dirtyPath, discardPrefix) {
		cleanPath = dirtyPath[len(discardPrefix):]
	}
	return cleanPath
}

func HasFilePath(filePath string) bool {
	f, err := CoreFs.Open(Clean(filePath))
	if err == nil {
		f.Close()
		return true
	}
	return false
}

func FindFile(fileName string) (filePaths []string, err error) {
	return []string{}, nil

}

// find a filename in the list of embedded file systems
func FindEmbeddedFile(filename string) ([]string, error) {
	slog.Debug(fmt.Sprintf(" searching for (%s) in core embedded file system", filename))

	matches := []string{}

	//start at the root folder
	return searchFolder(CoreFs, "", filename, matches)
}

// searchFolder walks a path looking for filename is the efs file system.
// It will append the path to that file in the matches slice.
// it only returns an error if there is an error walking the file system. Not
// found is indicated by matches not growing with the function is called.
func searchFolder(efs embed.FS, path string, filename string, matches []string) ([]string, error) {
	slog.Debug("search folder: " + path)
	found := false
	// walk this FS for any matches
	err := fs.WalkDir(efs, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Base(p) == filename {
			slog.Debug(">>found: " + p)
			matches = append(matches, p)
			found = true
			return nil
		}
		slog.Debug("   skip: " + p)
		return nil
	})

	slog.Debug("end of search for " + filename)

	if found {
		return matches, err
	} else if err == nil {
		// return a simple not found
		return matches, err
	}

	// no match in this folder check subfolders
	return matches, errors.New(err.Error() + " (" + filename + " not found)")
}
