package checker

import (
	"os"
	"testing"

	"github.com/cnrancher/hangar/pkg/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

func Test_annotationCheck(t *testing.T) {
	p, _ := utils.GetAbsPath("../../../")
	logrus.Infof("chart path: %v", p)
	checker := NewChecker(p, "v2.7")
	if err := checker.init(); err != nil {
		t.Error(err)
		return
	}
	err := checker.annotationCheck()
	if err != nil {
		t.Error(err)
	}
}

func Test_imageCheck(t *testing.T) {
	p, _ := utils.GetAbsPath("../../../")
	logrus.Infof("chart path: %v", p)
	checker := NewChecker(p, "v2.7")
	if err := checker.init(); err != nil {
		t.Error(err)
		return
	}
	err := checker.imageCheck()
	if err != nil {
		t.Error(err)
	}
}

func Test_systemDefaultRegistryCheck(t *testing.T) {
	p, _ := utils.GetAbsPath("../../../")
	logrus.Infof("chart path: %v", p)
	checker := NewChecker(p, "v2.7")
	if err := checker.init(); err != nil {
		t.Error(err)
		return
	}
	err := checker.systemDefaultRegistryCheck()
	if err != nil {
		t.Error(err)
	}
}
