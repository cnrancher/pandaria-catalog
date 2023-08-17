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
	for k, v := range pluginContainerImages {
		assert.Equal(t, len(v), len(testedContainerImages[k]),
			"some container image of plugin [%v] not tested, expted:\n%v\n actual: %v",
			k, printObject(v), printObject(testedContainerImages[k]))
	}
}
