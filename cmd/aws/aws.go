//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package aws adds management functions for Amazon Web Services

package aws

import (
	"log/slog"
	"runtime"

	"github.com/mrmxf/clog/ux/ui"
	"github.com/spf13/cobra"
)

// awsCmd represents the base aws command when called without any subcommands
var Command = &cobra.Command{
	Use:   "Aws",
	Short: "AWS utilities - dns etc",

	// Run interactively unless told to be batch / server
	Run: func(awsCmd *cobra.Command, args []string) {
		//run all sub commands interactively because no subcommand was requested
		_, err := ui.HomeMenu(awsCmd)
		if err != nil {
			slog.Error("cannot run AWS home menu" + err.Error())
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.AddCommand(DnsCommand)
}
