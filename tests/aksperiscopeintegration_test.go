package main

import (
	"fmt"
	"testing"

	"github.com/Azure/aks-periscope/pkg/utils"
)

func TestEndToEndIntegrationSuccessCase(t *testing.T) {
	runperiscopedeploycommand(t, false)
}

func TestEndToEndIntegrationUnsuccessFulCase(t *testing.T) {
	runperiscopedeploycommand(t, true)
}

func runperiscopedeploycommand(t *testing.T, validate bool) {
	// This flag switch on and off for storage account validation.
	validateflag := fmt.Sprintf("--validate=%v", validate)

	output, err := utils.RunCommandOnContainer("kubectl", "apply", "-f", "../deployment/aks-periscope.yaml", validateflag)

	if err != nil && validate {
		t.Logf("unsuccessful output: %v\n", err)
	}

	if output != "" && !validate {
		t.Logf("successful output: %v\n", output)
	}
}
