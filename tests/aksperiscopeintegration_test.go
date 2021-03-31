package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Azure/aks-periscope/pkg/utils"
	. "github.com/onsi/gomega"
)

func TestEndToEndIntegrationSuccessCase(t *testing.T) {
	runperiscopedeploycommand(t, false)
	g := NewGomegaWithT(t)

	// check if pods are running
	g.Eventually(func() bool {
		return checkifpodsrunning(t)
	}, "60s", "2s").Should(BeTrue())

	// check if location of the logs is not empty
	g.Eventually(func() bool {
		return islogsdirempty(t)
	}, "60s", "2s").Should(BeFalse())

}

func TestEndToEndIntegrationUnsuccessFulCase(t *testing.T) {
	runperiscopedeploycommand(t, true)
}

func runperiscopedeploycommand(t *testing.T, validate bool) {
	// This flag switch on and off for storage account validation.
	validateflag := fmt.Sprintf("--validate=%v", validate)
	g := NewGomegaWithT(t)

	output, err := utils.RunCommandOnContainer("kubectl", "apply", "-f", "../deployment/aks-periscope.yaml", validateflag)

	if err != nil && validate {
		g.Expect(err).Should(HaveOccurred())
		t.Logf("unsuccessful output: %v\n", err)
	}

	if output != "" && !validate {
		g.Expect(err).ToNot(HaveOccurred())
		t.Logf("successful output: %v\n", output)
	}
}

func checkifpodsrunning(t *testing.T) bool {
	g := NewGomegaWithT(t)

	output, err := utils.RunCommandOnContainer("kubectl", "get", "pods", "-n", "aks-periscope")
	firstpod := strings.Split(output, "\n")

	firstpodname := strings.Fields(firstpod[1])[0]
	firstpodstate := strings.Fields(firstpod[1])[2]

	t.Logf(" Outcome is %v ===> %v", firstpodname, firstpodstate)

	if firstpodstate == "Running" {
		return true
	}

	if err != nil {
		g.Expect(err).ToNot(HaveOccurred())
		t.Logf("unsuccessful error: %v\n", err)
	}

	if output != "" {
		g.Expect(err).ToNot(HaveOccurred())
		t.Logf("successful output: %v\n", output)
	}

	return false
}

func islogsdirempty(t *testing.T) bool {
	knownloglocation := "/var/log/aks-periscope/"

	f, err := os.Open(knownloglocation)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}

	return false // Either not empty or error, suits both cases
}
