package collector

import (
	"github.com/Azure/aks-periscope/pkg/interfaces"
	"github.com/Azure/aks-periscope/pkg/utils"
	"path/filepath"
)

// HelmCollector defines a Helm Collector struct
type HelmCollector struct {
	BaseCollector
}

var _ interfaces.Collector = &HelmCollector{}

// NewHelmCollector is a constructor
func NewHelmCollector(exporter interfaces.Exporter) *HelmCollector {
	return &HelmCollector{
		BaseCollector: BaseCollector{
			collectorType: Helm,
			exporter:      exporter,
		},
	}
}

// Collect implements the interface method
func (collector *HelmCollector) Collect() error {
	rootPath, err := utils.CreateCollectorDir(collector.GetName())
	if err != nil {
		return err
	}

	output, err := utils.RunCommandOnContainer("apk", "add", "curl", "openssl", "bash", "--no-cache")
	if err != nil {
		return err
	}

	output, err = utils.RunCommandOnContainer("curl", "-fsSl", "-o", "get_helm.sh", "https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3")
	if err != nil {
		return err
	}

	output, err = utils.RunCommandOnContainer("chmod", "+x", "get_helm.sh")
	if err != nil {
		return err
	}

	output, err = utils.RunCommandOnContainer("./get_helm.sh")
	if err != nil {
		return err
	}

	helmListFile := filepath.Join(rootPath, "helm_list")
	output, err = utils.RunCommandOnContainer("helm", "list", "--all-namespaces")
	if err != nil {
		errorMessage := err.Error()
		err = utils.WriteToFile(helmListFile, errorMessage)
		if err != nil {
			return err
		}
	}

	err = utils.WriteToFile(helmListFile, output)
	if err != nil {
		return err
	}

	collector.AddToCollectorFiles(helmListFile)

	helmHistoryFile := filepath.Join(rootPath, "helm_history")
	output, err = utils.RunCommandOnContainer("helm", "history", "-n", "default", "azure-arc")
	if err != nil {
		errorMessage := err.Error()
		err = utils.WriteToFile(helmHistoryFile, errorMessage)
		if err != nil {
			return err
		}
	}

	err = utils.WriteToFile(helmHistoryFile, output)
	if err != nil {
		return err
	}

	collector.AddToCollectorFiles(helmHistoryFile)

	return nil
}
