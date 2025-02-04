// Copyright ©2018-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package semver_test

import (
	"embed"
	"reflect"
	"testing"

	"github.com/mrmxf/clog/semver"
	. "github.com/smartystreets/goconvey/convey"
)

var refSemverInfo =""
var refLinkerDataDefault= "hash|date|suffix|app|title"
func refInitialise(fs embed.FS, filePath string) error {return nil}

// compareAllFields does a deep inspect of the properties
// values are of different types so reflect.DeepEqual will always be false
// forcing a manual check
func compareAllFields(t *testing.T, aStruct interface{}, bStruct interface{}) bool {
	a := reflect.ValueOf(aStruct)
	b := reflect.ValueOf(bStruct)

	for i := 0; i < a.NumField(); i++ {
		aVal := a.Field(i).Interface()
		bVal := b.Field(i).Interface()
		if aVal != bVal {
			// t.Log("aVal =", aVal, "bVal =", bVal)
			return false
		}
	}
	return true
}

func TestSpec(t *testing.T) {

// Check exported elements (backwards compatibility)
Convey("We should have consistent exported elements", t, func() {

	Convey("Exported types", func() {
		Convey("type LDgolangLinkerData should not have changed", func() {
			unchanged := compareAllFields(t, semver.LDgolangLinkerData{},refLinkerData{})
			So(unchanged, ShouldBeTrue)
		})
		Convey("type VersionInfo should not have changed", func() {
			unchanged := compareAllFields(t, semver.VersionInfo{},refVersionInfo{})
			So(unchanged, ShouldBeTrue)
		})
		Convey("type ReleaseHistory should not have changed", func() {
			unchanged := compareAllFields(t, semver.ReleaseHistory{},refReleaseHistory{})
			So(unchanged, ShouldBeTrue)
		})
	})
	Convey("Exported properties", func() {
		Convey("SemVerInfo should exist", func() {
			So(semver.SemVerInfo, ShouldEqual, refLinkerDataDefault)
		})
		Convey("SemVerInfo should be the right type", func() {
			So(semver.SemVerInfo, ShouldHaveSameTypeAs, refSemverInfo)
		})
		Convey("LinkerDataDefault should exist", func() {
			So(semver.LinkerDataDefault, ShouldEqual, refLinkerDataDefault)
		})
		Convey("LinkerDataDefault should be the right type", func() {
			So(semver.LinkerDataDefault, ShouldHaveSameTypeAs, refLinkerDataDefault)
		})
	})

	Convey("Exported functions", func() {
		// Initialise
		Convey("Initialise should exist", func() {
			So(semver.Initialise, ShouldNotBeNil)
		})
		Convey("UsePrettyLogger should be the right type", func() {
			So(semver.Initialise, ShouldHaveSameTypeAs, refInitialise)
		})


	})

})}
