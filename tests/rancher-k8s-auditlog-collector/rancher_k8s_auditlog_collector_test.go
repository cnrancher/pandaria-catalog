package rancher_k8s_auditlog_collector

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
