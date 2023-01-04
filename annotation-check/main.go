package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ChartYaml struct {
	Annotations ChartAnnotations `yaml:"annotations"`
	Name        string           `yaml:"name"`
	Version     string           `yaml:"version"`
}

type ChartAnnotations struct {
	KubeVersion    string `yaml:"catalog.cattle.io/kube-version"`
	RancherVersion string `yaml:"catalog.cattle.io/rancher-version"`
}

var (
	NoRancherVersionFile = "no-rancher-version.txt"
	NoKubeVersionFile    = "no-kube-version.txt"
)

func init() {
	// logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		logrus.Info("This program ensure charts in this repo has" +
			"'rancher-version' and 'kube-version' annotations")
		logrus.Info("Usage: go run ./main.go <path>")
		return
	}

	// delete output file if exists, ignore error
	os.Remove(NoKubeVersionFile)
	os.Remove(NoRancherVersionFile)

	path := os.Args[1]
	flag := false
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			logrus.Infof("failed to access %q: %v", path, err)
			return err
		}
		if info.IsDir() {
			if isSubChart(path) {
				logrus.Debugf("skip subchart folder: %s", path)
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() != "Chart.yaml" {
			return nil
		}
		logrus.Debugf("visit file: %q", path)
		if err := checkAnnotation(path); err != nil {
			logrus.Errorf("Failed to check %q", path)
			logrus.Error(err)
			flag = true
		}
		return nil
	})
	if err != nil {
		logrus.Fatalf("Error walking the path %q: %v\n", path, err)
	}
	if flag {
		logrus.Fatal("There are some charts failed to check")
	}
}

func checkAnnotation(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("checkAnnotation: %w", err)
	}

	spec := ChartYaml{}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("checkAnnotation: %w", err)
	}
	var errmsg string
	if spec.Annotations.KubeVersion == "" {
		errmsg += fmt.Sprintf("%q version %q does not have 'kube-version'\n",
			spec.Name, spec.Version)
		appendFileLine(NoKubeVersionFile,
			fmt.Sprintf("%s - %s", spec.Name, spec.Version))
	}
	if spec.Annotations.RancherVersion == "" {
		errmsg += fmt.Sprintf("%q version %q does not have 'rancher-version'\n",
			spec.Name, spec.Version)
		appendFileLine(NoRancherVersionFile,
			fmt.Sprintf("%s - %s", spec.Name, spec.Version))
	}
	if errmsg != "" {
		return errors.New(errmsg)
	}
	return nil
}

func isSubChart(path string) bool {
	spec := strings.Split(path, "/")
	var firstChart bool = false
	for _, v := range spec {
		if v == "packages" {
			// skip packages folder
			return true
		}
		if v == "charts" {
			if !firstChart {
				firstChart = true
			} else {
				return true
			}
		}
	}

	return false
}

func appendFileLine(fileName string, line string) error {
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
