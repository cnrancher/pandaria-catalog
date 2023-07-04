package chartpackages

import (
	"archive/tar"
	"fmt"
	"io"
	"path/filepath"

	"github.com/klauspost/pgzip"
	"helm.sh/helm/v3/pkg/chart"
	k8sYaml "sigs.k8s.io/yaml"
)

func LoadMetadataTgz(tgz io.Reader) (*chart.Metadata, error) {
	gzr, err := pgzip.NewReader(tgz)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil, fmt.Errorf("can not found Chart.yaml")
		case err != nil:
			return nil, err
		case header.Typeflag == tar.TypeReg:
			fileName := filepath.Base(header.Name)
			if fileName != "Chart.yaml" && fileName != "Chart.yml" {
				continue
			}
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			metadata := new(chart.Metadata)
			if err := k8sYaml.Unmarshal(data, metadata); err != nil {
				return metadata, fmt.Errorf("can not load Chart.yaml: %w", err)
			}
			return metadata, nil
		default:
			continue
		}
	}
}
