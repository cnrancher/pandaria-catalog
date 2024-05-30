package chartpackages

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/klauspost/pgzip"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chart"
	k8sYaml "sigs.k8s.io/yaml"
)

func LoadMetadataTgz(tgz io.Reader, patchFile string) (*chart.Metadata, error) {
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
			if patchFile != "" {
				tmpDir, err := os.MkdirTemp("", "chart-check-*")
				if err != nil {
					return nil, fmt.Errorf("mkdirTemp failed: %w", err)
				}
				tmpFile, err := os.Create(filepath.Join(tmpDir, fileName))
				if err != nil {
					return nil, fmt.Errorf("os.Create failed: %w", err)
				}
				_, err = tmpFile.Write(data)
				if err != nil {
					return nil, fmt.Errorf("write failed: %w", err)
				}

				tmpFile.Close()
				if err = ApplyPatch(patchFile, tmpDir); err != nil {
					return nil, fmt.Errorf("failed to apply patch: %w", err)
				}
				data, err = os.ReadFile(filepath.Join(tmpDir, fileName))
				if err != nil {
					return nil, fmt.Errorf("failed to read file: %w", err)
				}
				os.RemoveAll(tmpDir)
				logrus.Infof("Applied patch: %v", patchFile)
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

// ApplyPatch applies a patch file located at patchPath to the destDir on the filesystem
func ApplyPatch(patchPath, destDir string) error {
	// TODO: (aiyengar2): find a better library to actually generate and apply patches
	// There doesn't seem to be any existing library at the moment that can work with unified patches
	pathToPatchCmd, err := exec.LookPath("patch")
	if err != nil {
		return fmt.Errorf("cannot generate patch file if GNU patch is not available")
	}

	var buf bytes.Buffer
	patchFile, err := os.Open(patchPath)
	if err != nil {
		return err
	}
	defer patchFile.Close()

	cmd := exec.Command(pathToPatchCmd, "-E", "-p1")
	cmd.Dir = destDir
	cmd.Stdin = patchFile
	cmd.Stdout = &buf

	if err = cmd.Run(); err != nil {
		logrus.Errorf("\n%s", &buf)
		err = fmt.Errorf("unable to generate patch with error: %s", err)
	}
	return err
}
