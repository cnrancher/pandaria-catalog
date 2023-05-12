package main

import (
	"flag"
	"fmt"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	hangarUtils "github.com/cnrancher/hangar/pkg/utils"
	"github.com/cnrancher/pandaria-catalog/check/pkg/checker"
	"github.com/cnrancher/pandaria-catalog/check/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	cmdDebug       bool
	rancherVersion string
)

func main() {
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        false,
		TimestampFormat: "15:04:05", // hour, time, sec only
	})
	flag.BoolVar(&cmdDebug, "debug", false, "Enable the debug output")
	flag.StringVar(&rancherVersion, "version", "v2.7", "Rancher Version, v2.7 or v2.6")
	flag.Usage = func() {
		logrus.Infof("Usage:   ./check [OPTIONS] <path>")
		logrus.Infof("Example: ./check --debug --version=v2.7 ./charts")
		logrus.Infof("Available options:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if cmdDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if v, err := hangarUtils.EnsureSemverValid(rancherVersion); err != nil {
		logrus.Fatalf("invalid version: %v", rancherVersion)
	} else {
		rancherVersion = v
	}

	if len(flag.Args()) == 0 {
		logrus.Error("project dir not specified")
		flag.Usage()
		return
	}

	// Delete output files if exists
	os.Remove(utils.NoKubeVersionFile)
	os.Remove(utils.NoRancherVersionFile)
	os.Remove(utils.ImageCheckFailedFile)
	os.Remove(utils.SystemDefaultRegistryCheckFailed)

	checker := checker.NewChecker(flag.Args()[0], rancherVersion)
	if err := checker.Check(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("ALL CHECK PASSED")
}

func AppendFileLine(fileName string, line string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("AppendFileLine: %w", err)
	}
	if _, err := f.Write([]byte(line + "\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		return fmt.Errorf("AppendFileLine: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("AppendFileLine: %w", err)
	}

	return nil
}
