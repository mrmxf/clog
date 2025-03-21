// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package slogger_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"time"

	slog "github.com/mrmxf/clog/slogger"
	. "github.com/smartystreets/goconvey/convey"
)

func stripEscapes(s string) string {
	inEscSeq := false
	var res string = ""

	for _, r := range s {
		if inEscSeq {
			if r == 'm' {
				inEscSeq = false
			}
		} else {
			if r == '\x1b' {
				inEscSeq = true
			} else {
				res += string(r)
			}
		}
	}
	return res
}

func checkLogOutput(t *testing.T, buf bytes.Buffer, title string, expected string) {
	lenPrefix := len("2025-03-21 14:50:21 ")
	now := time.Time()
	bare := stripEscapes()
	Convey(fmt.Sprintf("%s check", title), func() {

		Convey("bufio len count", func() {
			So(buf.Len(), ShouldEqual, lenPrefix+len(expected))
		})
		Convey("output string == \"expected\"", func() {
			s := buf.String()
			So(stripEscapes(s), ShouldEqual, expected)
		})
		Convey("time stamp", func() {
			logT := time.Parse()
			So(stripEscapes(s), ShouldEqual, expected)
		})

	})
}

func TestSpec_Levels(t *testing.T) {
	// Check each level works correctly
	Convey("Each level should work properly", t, func() {
		var buf bytes.Buffer
		out := bufio.NewWriter(&buf)
		slog.UsePrettyIoLogger(out, slog.LevelDebug)

		out.Reset(&buf)
		slog.Trace("Trace")
		out.Flush()
		checkLogOutput(t, buf, "Trace", "--- Trace")

		out.Reset(&buf)
		slog.Debug("Debug")
		out.Flush()
		checkLogOutput(t, buf, "Debug", "DBG Debug")

		out.Reset(&buf)
		slog.Info("Info")
		out.Flush()
		checkLogOutput(t, buf, "Info", "INF Info")

		out.Reset(&buf)
		slog.Success("Success")
		out.Flush()
		checkLogOutput(t, buf, "Success", " OK Success")

		out.Reset(&buf)
		slog.Warn("Warn")
		out.Flush()
		checkLogOutput(t, buf, "Warn", "WRN Warn")

		out.Reset(&buf)
		slog.Error("Error")
		out.Flush()
		checkLogOutput(t, buf, "Error", "ERR Error")

		out.Reset(&buf)
		slog.Fatal("Fatal")
		out.Flush()
		checkLogOutput(t, buf, "Fatal", "FTL Fatal")

		out.Reset(&buf)
		slog.Emergency("Emergency")
		out.Flush()
		checkLogOutput(t, buf, "Emergency", "!!! Emergency")

	})

}
