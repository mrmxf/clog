//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Ls
// returns:
//   ..a list of files in the embedded fs

package list

import (
	"fmt"
	"io/fs"
	"log/slog"
	"runtime"

	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/crayon"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "List",
	Short: "list all embedded files to  stdout",
	Long:  `errors logged to stderr with status 1.`,

	Run: func(cmd *cobra.Command, args []string) {
		clogFs := config.CoreFs()
		c := crayon.Color()

		// List the files to update
		err := fs.WalkDir(clogFs, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				slog.Error(err.Error())
			}
			clogFsStat, _ := fs.Stat(clogFs, path)

			if clogFsStat.IsDir() {
				fmt.Printf("%v/\n", c.I(path))
			} else {
				fmt.Printf("%v\n", path)
			}
			return nil
		})
		if err != nil {
			slog.Error("Cannot List embedded files" + err.Error())
		}

	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
