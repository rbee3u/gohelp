package conflog_test

import (
	"testing"

	"github.com/rbee3u/gohelp/conflog"
)

func TestLogger(t *testing.T) {
	l, _ := conflog.New(conflog.WithReportCaller("enable"))
	l.Info("Hello")
}
