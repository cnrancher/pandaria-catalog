package rancher_gpu_sharing

import (
	"testing"

	"github.com/rancher/hull/pkg/test"
)

func TestChart(t *testing.T) {
	opts := test.GetRancherOptions()
	opts.HelmLint.Rancher.Enabled = false
	opts.Coverage.Disabled = true
	// opts.Coverage.IncludeSubcharts = false
	// opts.YAMLLint.Enabled = false
	suite.Run(t, opts)
}
