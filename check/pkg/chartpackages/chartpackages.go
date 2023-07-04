package chartpackages

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cnrancher/hangar/pkg/rancher/chartimages"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
)

// Package represents the configuration of a particular forked Helm chart
type Package struct {
	URL            string `yaml:"url,omitempty"`
	Name           string `yaml:"name,omitempty"`
	Version        string `yaml:"version,omitempty"`
	PackageVersion *int   `yaml:"packageVersion,omitempty"`
	DoNotRelease   bool   `yaml:"doNotRelease,omitempty"`
	WorkingDir     string `yaml:"workingDir,omitempty"`
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
		if _, err := os.Stat(filepath.Join(p, "package.yaml")); err != nil {
			return nil
		}
		pkg, err := decodePackageYaml(filepath.Join(p, "package.yaml"))
		if err != nil {
			logrus.Warnf("decodePackageYaml: %v", err)
			return nil
		}
		if skipPackage(pkg) {
			logrus.Debugf("chart packages skip pkg %v", p)
			return filepath.SkipDir
		}

		var metadata *chart.Metadata
		if pkg.URL == "" || pkg.URL == "local" {
			metadata, err = chartimages.LoadMetadata(filepath.Join(p, pkg.WorkingDir))
			if err != nil {
				return err
			}
		} else {
			client := http.Client{
				Timeout: time.Second * 5,
			}
			resp, err := client.Get(pkg.URL)
			if err != nil {
				return fmt.Errorf("failed to get %q: %w", pkg.URL, err)
			}
			defer resp.Body.Close()
			metadata, err = LoadMetadataTgz(resp.Body)
			if err != nil {
				return err
			}
		}
		if metadata == nil {
			logrus.Warnf("failed to get chart metadata: %+v", pkg)
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
