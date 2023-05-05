package checker

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	hangarUtils "github.com/cnrancher/hangar/pkg/utils"
	"github.com/klauspost/pgzip"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/cnrancher/hangar/pkg/rancher/chartimages"
	"github.com/cnrancher/pandaria-catalog/check/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Checker struct {
	Path string

	index *repo.IndexFile
}

func NewChecker(p string) *Checker {
	return &Checker{
		Path: p,
	}
}

func (cc *Checker) init() error {
	var err error
	cc.index, err = repo.LoadIndexFile(
		filepath.Join(cc.Path, "index.yaml"))
	if err != nil {
		return fmt.Errorf("LoadIndexFile: %q: %v", cc.Path, err)
	}
	if len(cc.index.Entries) == 0 {
		return fmt.Errorf("no charts found in %q", cc.Path)
	}
	return nil
}

func (cc *Checker) Check() error {
	if cc.Path == "" {
		return fmt.Errorf("chart path not specified")
	}
	var err error
	if cc.Path, err = hangarUtils.GetAbsPath(cc.Path); err != nil {
		return err
	}

	if err := cc.init(); err != nil {
		return err
	}
	var failed bool
	if err := cc.annotationCheck(); err != nil {
		logrus.Error(err)
		failed = true
	}
	if err := cc.imageCheck(); err != nil {
		logrus.Error(err)
		failed = true
	}
	if err := cc.systemDefaultRegistryCheck(); err != nil {
		logrus.Error(err)
		failed = true
	}

	if failed {
		return fmt.Errorf("check failed")
	}
	return nil
}

func (cc *Checker) annotationCheck() error {
	logrus.Infof("Start annotation check...")
	rvList := []string{}
	kvList := []string{}
	annotationCheckFailed := false
	for _, versions := range cc.index.Entries {
		if len(versions) == 0 {
			continue
		}
		lv := versions[0]
		rv, ok := lv.Annotations[utils.RancherVersionAnnotationKey]
		if !ok {
			nr := fmt.Sprintf("%s - %s", lv.Name, lv.Version)
			logrus.Errorf("FAILED: No rancher-version annotation: %s", nr)
			rvList = append(rvList, nr)
		} else {
			logrus.Infof("Found rancher-version of %q: %q",
				lv.Name, rv)
		}
		kv, ok := lv.Annotations[utils.KubeVersionAnnotationKey]
		if !ok {
			nk := fmt.Sprintf("%s - %s", lv.Name, lv.Version)
			logrus.Errorf("FAILED: No kube-version annotation: %s", nk)
			kvList = append(kvList, nk)
		} else {
			logrus.Infof("Found kube-version of %q: %q",
				lv.Name, kv)
		}
	}

	if len(rvList) != 0 {
		hangarUtils.SaveSlice(utils.NoRancherVersionFile, rvList)
		annotationCheckFailed = true
	}
	if len(kvList) != 0 {
		hangarUtils.SaveSlice(utils.NoKubeVersionFile, kvList)
		annotationCheckFailed = true
	}
	if annotationCheckFailed {
		return fmt.Errorf("annotation check failed")
	} else {
		logrus.Infof("Annotation check passed")
	}

	return nil
}

func (cc *Checker) imageCheck() error {
	logrus.Infof("Start image check...")
	failedList := []string{}
	crdRegexp, err := regexp.Compile(".*-crd")
	if err != nil {
		logrus.Warnf("regexp.Compile: %v", err)
	}

	for _, versions := range cc.index.Entries {
		if len(versions) == 0 {
			continue
		}
		lv := versions[0]
		if crdRegexp != nil && crdRegexp.MatchString(lv.Name) {
			logrus.Infof("Skip check CRD chart: %v - %v", lv.Name, lv.Version)
			continue
		}
		if len(lv.URLs) == 0 {
			logrus.Warnf("URL of %q is empty", lv.Name)
			continue
		}
		values, err := cc.getChartValues(lv.URLs[0])
		if err != nil {
			logrus.Error(err)
			continue
		}
		var imageSet = make(map[string]map[string]bool)
		err = chartimages.PickImagesFromValuesMap(
			imageSet, values[0], "", chartimages.Linux)
		if err != nil {
			logrus.Warn(err)
		}
		err = chartimages.PickImagesFromValuesMap(
			imageSet, values[0], "", chartimages.Windows)
		if err != nil {
			logrus.Warn(err)
		}
		var foundImage bool = false
		for range imageSet {
			foundImage = true
			break
		}
		msg := fmt.Sprintf("%s - %s", lv.Name, lv.Version)
		if foundImage {
			logrus.Infof("PASS: %v", msg)
		} else {
			logrus.Errorf("FAILED: no images found from chart: %v", msg)
			failedList = append(failedList, msg)
		}
	}
	if len(failedList) != 0 {
		hangarUtils.SaveSlice(utils.NoKubeVersionFile, failedList)
		return fmt.Errorf("chart image check failed")
	}
	logrus.Infof("Image check passed")

	return nil
}

func (cc *Checker) systemDefaultRegistryCheck() error {
	logrus.Infof("Start systemDefaultRegistry check...")
	failedList := []string{}
	crdRegexp, err := regexp.Compile(".*-crd")
	if err != nil {
		logrus.Warnf("regexp.Compile: %v", err)
	}
	for _, versions := range cc.index.Entries {
		if len(versions) == 0 {
			continue
		}
		lv := versions[0]
		if crdRegexp != nil && crdRegexp.MatchString(lv.Name) {
			logrus.Infof("Skip check CRD chart: %v - %v", lv.Name, lv.Version)
			continue
		}
		msg := fmt.Sprintf("%s: %s", lv.Name, lv.URLs[0])
		values, err := cc.getChartValues(lv.URLs[0])
		if err != nil {
			logrus.Errorf("failed to get %q values.yaml: %v", lv.URLs[0], err)
			failedList = append(failedList, msg)
			continue
		}
		if len(values) == 0 {
			logrus.Warnf("%q does not have values.yaml", msg)
			continue
		}
		tmplString, err := cc.getChartHelperTemplate(lv.URLs[0])
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				logrus.Warnf("%q does not have '_helpers.tpl'", lv.URLs[0])
				// failedList = append(failedList, msg)
				continue
			}
			logrus.Errorf("failed to get the content of %q '_helpers.tpl': %v",
				lv.URLs[0], err)
			failedList = append(failedList, msg)
			continue
		}

		// Add fake template functions here
		var funcMap = template.FuncMap{
			"default":    fakeFunction,
			"include":    fakeFunction,
			"trunc":      fakeFunction,
			"trimSuffix": fakeFunction,
			"contains":   fakeFunction,
			"replace":    fakeFunction,
			"quote":      fakeFunction,
			"toYaml":     fakeFunction,
			"fromYaml":   fakeFunction,
			"empty":      fakeFunction,
			"has":        fakeFunction,
			"hasKey":     fakeFunction,
			"keys":       fakeFunction,
			"deepCopy":   fakeFunction,
			"set":        fakeFunction,
			"sha256sum":  fakeFunction,
			"list":       fakeFunction,
		}
		tmplString = fmt.Sprintf("%v\nTEST: {{ template %q . }}\n",
			tmplString, "system_default_registry")
		tmpl, err := template.New(lv.URLs[0]).
			Funcs(funcMap).Parse(tmplString)
		if err != nil {
			logrus.Errorf("failed to parse template %q: %v", lv.URLs[0], err)
			failedList = append(failedList, msg)
			continue
		}

		// Check global.systemDefaultRegistry or
		// global.cattle.systemDefaultRegistry is defined in values.yaml
		global, ok := values[0]["global"].(map[any]any)
		if !ok {
			values[0]["global"] = make(map[any]any)
			global = values[0]["global"].(map[any]any)
		}
		global["systemDefaultRegistry"] = "registry.example.io"
		cattle, ok := global["cattle"].(map[any]any)
		if !ok {
			global["cattle"] = make(map[any]any)
			cattle = global["cattle"].(map[any]any)
		}
		cattle["systemDefaultRegistry"] = "registry.example.io"

		buff := bytes.Buffer{}
		err = tmpl.Execute(&buff, struct {
			Values any `json:"Values"`
		}{Values: values[0]})
		if err != nil {
			logrus.Errorf("FAILED: template.Execute failed: %v", err)
			failedList = append(failedList, msg)
			continue
		}
		rendered := buff.String()

		if strings.Contains(rendered, "TEST: registry.example.io") {
			// the systemDefaultRegistry is effective
			logrus.Infof("PASS: %v - %v", lv.Name, lv.Version)
		} else {
			logrus.Errorf("FAILED: %s", msg)
			failedList = append(failedList, msg)
		}
	}

	if len(failedList) != 0 {
		hangarUtils.SaveSlice(
			utils.SystemDefaultRegistryCheckFailed, failedList)
		return fmt.Errorf("systemDefaultRegistry check failed")
	}
	logrus.Infof("systemDefaultRegistry check passed")
	return nil
}

func (cc *Checker) getChartValues(chartURL string) ([]map[any]any, error) {
	chartPath := filepath.Join(cc.Path, chartURL)
	info, err := os.Stat(chartPath)
	if err != nil {
		logrus.Warnf("%q: %s", chartPath, err)
		return nil, err
	}

	if info.IsDir() {
		return chartimages.DecodeValuesInDir(chartPath)
	}
	return chartimages.DecodeValuesInTgz(chartPath)
}

func (cc *Checker) getChartHelperTemplate(chartURL string) (string, error) {
	info, err := os.Stat(filepath.Join(cc.Path, chartURL))
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		// destination file is a directory
		helpersTplPath := filepath.Join(
			cc.Path, chartURL, "templates", "_helpers.tpl")
		data, err := os.ReadFile(helpersTplPath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// destination file is a tarball
	tgz, err := os.Open(filepath.Join(cc.Path, chartURL))
	if err != nil {
		return "", err
	}
	defer tgz.Close()
	gzr, err := pgzip.NewReader(tgz)
	if err != nil {
		return "", err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	var data []byte
	// var headerName string
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return "", nil
		case err != nil:
			return "", err
		case header.Typeflag == tar.TypeReg &&
			filepath.Base(header.Name) == "_helpers.tpl":
			data, err = io.ReadAll(tr)
			if err != nil {
				return "", err
			}
			logrus.Debugf("Read: %v", header.Name)
			return string(data), nil
		default:
			continue
		}
	}
}

func fakeFunction() string {
	return ""
}
