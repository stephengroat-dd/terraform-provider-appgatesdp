package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApplianceCustomizationBasic(t *testing.T) {
	resourceName := "appgatesdp_appliance_customization.test_acc_appliance_customization"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceCustomizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceCustomizationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceCustomizationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "checksum_sha256", "e3a9fb24832dff49ea59ff79cff9b1f24cbc0974ec62ec700165a0631fee779e"),
					resource.TestCheckResourceAttr(resourceName, "file", "test-fixtures/appliance_customization_file.zip"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "size", "574"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccApplianceCustomizationImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckApplianceCustomizationBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_appliance_customization" "test_acc_appliance_customization" {
    name = "%s"
    file = "test-fixtures/appliance_customization_file.zip"

    tags = [
      "terraform",
      "api-created"
    ]
}
`, rName)
}

func testAccCheckApplianceCustomizationExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ApplianceCustomizationsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.ApplianceCustomizationsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching appliance customization with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckApplianceCustomizationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_appliance_customization" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ApplianceCustomizationsApi

		_, _, err := api.ApplianceCustomizationsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Appliance customization still exists, %+v", err)
		}
	}
	return nil
}

func testAccApplianceCustomizationImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
