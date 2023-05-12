package chartpackages

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cnrancher/hangar/pkg/rancher/chartimages"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/repo"
)

// Package represents the configuration of a particular forked Helm chart
type Package struct {
	Name           string `yaml:"name"`
	Version        string `yaml:"version,omitempty"`
	PackageVersion *int   `yaml:"packageVersion"`
	DoNotRelease   bool   `yaml:"doNotRelease,omitempty"`
}

func BuildIndex(dir string) (*repo.IndexFile, error) {
	var builtIndex = repo.NewIndexFile()

	err := filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			logrus.Warnf("%q: %v", p, err)
			return nil
		}
		if !i.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(p, "package.yaml")); err == nil {
			pkg, err := decodePackageYaml(filepath.Join(p, "package.yaml"))
			if err != nil {
				logrus.Warnf("decodePackageYaml: %v", err)
				return nil
			}
			if skipPackage(pkg) {
				logrus.Debugf("chart packages skip pkg %v", p)
				return filepath.SkipDir
			}
		}

		metadata, err := chartimages.LoadMetadata(p)
		if err != nil {
			return err
		}
		if metadata == nil {
			return nil
		}

		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return fmt.Errorf("building path for chart at %s: %w", dir, err)
		}

		err = builtIndex.MustAdd(metadata, rel, "", "")
		if err != nil {
			logrus.Warnf("failed to add %q into index file: %v", rel, err)
		}
		logrus.Debugf("chart packages add %v", p)
		return filepath.SkipDir
	})
	if err != nil {
		return nil, err
	}

	// sort index versions in descending order.
	builtIndex.SortEntries()

	return builtIndex, nil
}

func decodePackageYaml(name string) (*Package, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}
	pkg := &Package{}
	yaml.Unmarshal(data, pkg)

	return pkg, nil
}

func skipPackage(pkg *Package) bool {
	return pkg.DoNotRelease
}
