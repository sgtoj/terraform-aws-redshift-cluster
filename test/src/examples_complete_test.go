package test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Test the Terraform module in examples/complete using Terratest.
func TestExamplesComplete(t *testing.T) {
	t.Parallel()

  rand.Seed(time.Now().UnixNano())

  randId := strconv.Itoa(rand.Intn(100000))
  attributes := []string{randId}

  terraformOptions := &terraform.Options{
    // The path to where our Terraform code is located
    TerraformDir: "../../examples/complete",
    Upgrade:      true,
    // Variables to pass to our Terraform code using -var-file options
    VarFiles: []string{"fixtures.us-east-2.tfvars"},
    Vars: map[string]interface{}{
      "attributes": attributes,
    },
  }

  // At the end of the test, run `terraform destroy` to clean up any resources that were created
  defer terraform.Destroy(t, terraformOptions)

  // This will run `terraform init` and `terraform apply` and fail the test if there are any errors
  terraform.InitAndApply(t, terraformOptions)

  // Run `terraform output` to get the value of an output variable
  vpcCidr := terraform.Output(t, terraformOptions, "vpc_cidr")
  // Verify we're getting back the outputs we expect
  assert.Equal(t, "172.19.0.0/16", vpcCidr)

  // Run `terraform output` to get the value of an output variable
  privateSubnetCidrs := terraform.OutputList(t, terraformOptions, "private_subnet_cidrs")
  // Verify we're getting back the outputs we expect
  assert.Equal(t, []string{"172.19.0.0/19", "172.19.32.0/19"}, privateSubnetCidrs)

  // Run `terraform output` to get the value of an output variable
  publicSubnetCidrs := terraform.OutputList(t, terraformOptions, "public_subnet_cidrs")
  // Verify we're getting back the outputs we expect
  assert.Equal(t, []string{"172.19.96.0/19", "172.19.128.0/19"}, publicSubnetCidrs)

  // Run `terraform output` to get the value of an output variable
  clusterIdentifier := terraform.Output(t, terraformOptions, "id")
  expectedClusterIdentifier := "eg-test-redshift-cluster-" + randId
  // Verify we're getting back the outputs we expect
  assert.Equal(t, expectedClusterIdentifier, clusterIdentifier)

  // Run `terraform output` to get the value of an output variable
  arn := terraform.Output(t, terraformOptions, "arn")
  // Verify we're getting back the outputs we expect
  assert.Contains(t, arn, ":cluster:eg-test-redshift-cluster")

}
