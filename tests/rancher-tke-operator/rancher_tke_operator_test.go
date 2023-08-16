package rancher_tke_operator

import (
	"github.com/rancher/hull/pkg/test"
	"testing"
)

func TestChart(t *testing.T) {
	opts := test.GetRancherOptions()
	opts.Coverage.IncludeSubcharts = false
	opts.Coverage.Disabled = false
	opts.YAMLLint.Enabled = false
	suite.Run(t, opts)
}
