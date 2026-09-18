package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestMain(m *testing.M) {
	m.Run()
}

func GetTerraformOptions(t *testing.T) *terraform.Options {
	t.Helper()
	return &terraform.Options{
		TerraformDir: "../tf",
		Vars: map[string]interface{}{
			"project": "test-project",
		},
	}
}
