package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cnrancher/hangar/pkg/rancher/chartimages"
	"github.com/sirupsen/logrus"
)

const (
	RancherVersionAnnotationKey = "catalog.cattle.io/rancher-version"
	KubeVersionAnnotationKey    = "catalog.cattle.io/kube-version"
	NoRancherVersionFile        = "no-rancher-version.txt"
	NoKubeVersionFile           = "no-kube-version.txt"
)

var (
	cmdDebug bool
)

func main() {
	flag.BoolVar(&cmdDebug, "debug", false, "Enable the debug output")
	flag.Usage = func() {
		logrus.Infof("This program ensure the latest version chart in this " +
			"repo has 'rancher-version' and 'kube-version' annotations")
		logrus.Infof("Usage: annotation-check [OPTIONS] <path>")
		logrus.Infof("Available options:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if cmdDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	// Delete output files if exists
	os.Remove(NoKubeVersionFile)
	os.Remove(NoRancherVersionFile)

	path := args[0]
	var annotationCheckFailed = false
	index, err := chartimages.BuildOrGetIndex(path)
	if err != nil {
		logrus.Fatalf("Error walking the path %q: %v\n", path, err)
	}
	if len(index.Entries) == 0 {
		logrus.Warnf("No charts found in %q", path)
		return
	}

	for _, versions := range index.Entries {
		if len(versions) == 0 {
			continue
		}
		latestVersion := versions[0]
		// check annotations
		rv, ok := latestVersion.Annotations[RancherVersionAnnotationKey]
		if !ok {
			nr := fmt.Sprintf("%s - %s",
				latestVersion.Name, latestVersion.Version)
			logrus.Errorf("FAILED: No rancher-version annotation: %s", nr)
			annotationCheckFailed = true
			AppendFileLine(NoRancherVersionFile, nr)
		} else {
			logrus.Debugf("Found rancher-version of %q: %q",
				latestVersion.Name, rv)
		}
		kv, ok := latestVersion.Annotations[KubeVersionAnnotationKey]
		if !ok {
			nk := fmt.Sprintf("%s - %s",
				latestVersion.Name, latestVersion.Version)
			logrus.Errorf("FAILED: No kube-version annotation: %s", nk)
			annotationCheckFailed = true
			AppendFileLine(NoKubeVersionFile, nk)
		} else {
			logrus.Debugf("Found rancher-version of %q: %q",
				latestVersion.Name, kv)
		}
	}

	if annotationCheckFailed {
		logrus.Fatal("There are some charts failed to check")
	} else {
		logrus.Infof("annotation-heck passed")
	}
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
