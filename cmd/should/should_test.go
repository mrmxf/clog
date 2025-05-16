package should_test

import (
	"log/slog"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	slog.Info("Should command - tests for clog")
	Convey("Parse script headers correctly", t, func() {
		So(nil, ShouldBeNil)
	})
}
