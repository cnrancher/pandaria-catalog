package main

import (
	"fmt"
	"os"

	"github.com/cnrancher/hangar/pkg/rancher/chartimages"
	"github.com/sirupsen/logrus"
)

const (
	RancherVersionAnnotationKey = "catalog.cattle.io/rancher-version"
	KubeVersionAnnotationKey    = "catalog.cattle.io/kube-version"
)

var (
	NoRancherVersionFile = "no-rancher-version.txt"
	NoKubeVersionFile    = "no-kube-version.txt"
)

func init() {
	// logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		logrus.Infof("This program ensure charts in this repo has" +
			"'rancher-version' and 'kube-version' annotations")
		logrus.Infof("Usage: %s <path>", os.Args[0])
		return
	}

	// delete output file if exists, ignore error
	os.Remove(NoKubeVersionFile)
	os.Remove(NoRancherVersionFile)

	path := os.Args[1]
	var annotationCheckFailed = false
	index, err := chartimages.BuildOrGetIndex(path)
	if err != nil {
		logrus.Fatalf("Error walking the path %q: %v\n", path, err)
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
		}
		logrus.Debugf("Found rancher-version of %q: %q",
			latestVersion.Name, rv)
		kv, ok := latestVersion.Annotations[KubeVersionAnnotationKey]
		if !ok {
			nk := fmt.Sprintf("%s - %s",
				latestVersion.Name, latestVersion.Version)
			logrus.Errorf("FAILED: No kube-version annotation: %s", nk)
			annotationCheckFailed = true
			AppendFileLine(NoKubeVersionFile, nk)
		}
		logrus.Debugf("Found rancher-version of %q: %q",
			latestVersion.Name, kv)
	}

	if annotationCheckFailed {
		logrus.Fatal("There are some charts failed to check")
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
