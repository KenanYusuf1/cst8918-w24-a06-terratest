package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "b0004d43-233e-4db4-bf5c-5d1b0dbb3d8d"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix":"yusu0033",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}

func TestAzureVMNICExists(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "yusu0033",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initialize Terraform and apply configuration
	terraform.InitAndApply(t, terraformOptions)

	// Get the VM name and resource group from Terraform outputs
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Get the list of NICs attached to the VM
	nics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)

	// Assert that the VM has at least one NIC attached
	assert.NotEmpty(t, nics, "VM should have at least one NIC attached")
}

// TestUbuntuVersionOnVM checks if the deployed VM is running the expected Ubuntu version
func TestUbuntuVersionOnVM(t *testing.T) {
    terraformOptions := &terraform.Options{
        // The path to where our Terraform code is located
        TerraformDir: "../",
        // Override the default terraform variables
        Vars: map[string]interface{}{
            "labelPrefix": "yusu0033",
        },
    }

    // Ensure the Terraform is destroyed after test execution
    defer terraform.Destroy(t, terraformOptions)

    // Initialize and apply the Terraform code
    terraform.InitAndApply(t, terraformOptions)

    // Retrieve VM name and resource group from Terraform output
    vmName := terraform.Output(t, terraformOptions, "vm_name")
    resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

    // Retrieve the VM image details
    vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)

    // Assert the VM is running the expected Ubuntu version
    assert.Equal(t, "Canonical", vmImage.Publisher, "VM Publisher should be Canonical")
    assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer, "VM Offer should match Ubuntu Jammy")
    assert.Equal(t, "22_04-lts-gen2", vmImage.SKU, "VM SKU should match Ubuntu 22.04 LTS Gen2")
}
