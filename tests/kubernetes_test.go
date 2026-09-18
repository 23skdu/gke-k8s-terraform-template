package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestNamespaceExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	nsResource := plan.ResourcePlannedValuesMap["kubernetes_namespace_v1.prod"]
	require.NotNil(t, nsResource, "kubernetes_namespace_v1.prod should exist")
}

func TestNamespaceName(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	nsResource := plan.ResourcePlannedValuesMap["kubernetes_namespace_v1.prod"]
	require.NotNil(t, nsResource, "kubernetes_namespace_v1.prod should exist")

	metadata := nsResource.AttributeValues["metadata"]
	require.NotNil(t, metadata, "metadata should exist")
}
