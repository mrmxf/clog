//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package log adds a log command to the clog command line tool

package aws

import (
	"log/slog"
	"runtime"

	awsq "github.com/mrmxf/clog/aws"
	"github.com/mrmxf/clog/slogger"
	"github.com/spf13/cobra"
)

var debug bool
var byRegion bool
var byIp bool

// Command define the cobra settings for this command
var DnsCommand = &cobra.Command{
	Use:   "Dns",
	Short: "list DNS records for this AWS account sorted by DNS record name",
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			slogger.UsePrettyLogger(slog.LevelDebug)
		}

		// Configure sorting options according to your needs
		// This example sorts primarily by domain name, then by location, then by record value
		opts := awsq.SortOptions{
			Sort1: "name",     // Primary sort: case-insensitive domain name
			Sort2: "address",  // Secondary sort: AWS region/availability zone
			Sort3: "location", // Tertiary sort: DNS record value
		}

		if byIp {
			opts = awsq.SortOptions{
				Sort1: "address",  // Primary sort: case-insensitive domain name
				Sort2: "name",     // Secondary sort: AWS region/availability zone
				Sort3: "location", // Tertiary sort: DNS record value
			}
		}
		if byRegion {
			opts = awsq.SortOptions{
				Sort1: "location", // Primary sort: case-insensitive domain name
				Sort2: "name",     // Secondary sort: AWS region/availability zone
				Sort3: "address",  // Tertiary sort: DNS record value
			}
		}

		// Query all DNS records and print them to stdout
		// The function handles authentication via AWS credential chain
		err := awsq.Route53Domains(opts)
		if err != nil {
			slog.Error("Failed to query AWS domains", "err", err)
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&byRegion, "region", "R", false, "clog Aws Dns --region  ## sort by region first")
	Command.PersistentFlags().BoolVarP(&byIp, "ip", "I", false, "clog Aws Dns --debug  ## sort by IP first")
	Command.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "clog Aws Dns --debug  ## more logging")
}
