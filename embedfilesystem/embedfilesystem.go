//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/
// This file is part of clog.

// Package embedfilesystem manages the embedded files of clog and it works in
// conjunction with the `kfg` konfiguration package. If you are forking clog to
// make your own cli tool, then you probably want to override this package to
// include ONLY the files that you want in your own fork.

package embedfilesystem

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
)

//variable CoreFs is the exported file system

//go:embed all:deb
//go:embed all:tpl
//go:embed konfig.yaml
var CoreFs embed.FS

func HasFilePath(filePath string) bool {
	f, err := CoreFs.Open(filePath)
	if err == nil {
		f.Close()
		return true
	}
	return false
}

// find a filename in the list of embedded file systems
func FindEmbeddedFile(filename string) ([]string, error) {
	slog.Debug(fmt.Sprintf(" searching for (%s) in embedded file system", filename))

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
