//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/
// This file is part of clog.

package testclog

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func PrecedenceTest(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("precedence of script-snippet-internal", t, func() {
		x := 1

		Convey("snippets override internal", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})

		Convey("scripts override internal", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})

		Convey("scripts override snippets", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})

		Convey("scripts override snippets override internal", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
