//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// clog cfg find an embedded file in the list of File Systems
//
// usage:
// filePath:= cfg.FindEmbedded  ( "filename.ext" )
// filePath:= cfg.FindEmbeddedRe( "regexstr" )
//
// algorithm:
//   1. start each FS in order from 0 to max
//   2.   sort FS.folder contents alphabetically
//   3.   if match in FS.folder return path
//   4.   each subfolder search using step 2, match=>return else next
//   5.   next FS

package cfg

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
)

// find a filename in the list of embedded file systems
func FindEmbedded(filename string) (*embed.FS, []string, error) {
	slog.Debug(fmt.Sprintf(" searching for (%s) in %d file systems", filename, len(fsCache)))

	// iterate through the File Systems in lexical order
	for i, fs := range fsCache {
		matches := []string{}
		slog.Debug(fmt.Sprintf("---------- fs: %d of %d", i, len(fsCache)))

		//start at the root folder
		path := ""
		matches, err := searchFolder(fs, path, filename, matches)
		if err == nil {
			if len(matches) > 0 {
				slog.Debug(fmt.Sprintf(" Found %s in fs: %d", filename, i))
				return &fs, matches, err
			}
			// no error but no match - try next fs
			slog.Debug("no " + filename + " in fs: " + strconv.Itoa(i))
			continue
		}
	}
	msg := fmt.Sprintf("no %s in any fs", filename)
	slog.Debug(msg)
	return nil, []string{""}, errors.New(msg)
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

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
