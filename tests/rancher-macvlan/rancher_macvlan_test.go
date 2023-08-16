package rancher_macvlan

import (
	"testing"

	"github.com/rancher/hull/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestChart(t *testing.T) {
	opts := test.GetRancherOptions()
	opts.Coverage.IncludeSubcharts = false
	opts.Coverage.Disabled = false
	opts.YAMLLint.Enabled = false
	suite.Run(t, opts)
	assert.Equal(t, len(containerImages), len(testedContainerImages),
		"some container image not tested, expted: %v, actual: %v",
		containerImages, testedContainerImages)
}
