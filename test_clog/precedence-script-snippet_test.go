package test_clog

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func PrecedenceTests(t *testing.T) {

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
