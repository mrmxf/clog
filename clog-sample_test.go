package main_test

import (
	"log/slog"
	"testing"

	"github.com/mrmxf/clog/testclog"
)

func TestSpec(t *testing.T) {
	slog.Info("System tests for clog")
	t.Run("Precedence", testclog.PrecedenceTest)
}
