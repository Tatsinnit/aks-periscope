package collector

import (
	"github.com/Azure/aks-periscope/pkg/utils"
)

// DNSCollector defines a DNS Collector struct
type DNSCollector struct {
	data map[string]string
}

// NewDNSCollector is a constructor
func NewDNSCollector() *DNSCollector {
	return &DNSCollector{
		data: make(map[string]string),
	}
}

func (collector *DNSCollector) GetName() string {
	return "dns"
}

// Collect implements the interface method
func (collector *DNSCollector) Collect() error {
	result := make(map[string]string)

	output, err := utils.RunCommandOnHost("cat", "/etc/resolv.conf")
	if err != nil {
		return err
	}

	result["virtualmachine"] = output

	output, err = utils.RunCommandOnContainer("cat", "/etc/resolv.conf")
	if err != nil {
		return err
	}

	collector.data["kubernetes"] = output

	return nil
}

func (collector *DNSCollector) GetData() map[string]string {
	return collector.data
}
