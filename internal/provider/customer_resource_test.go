package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomerResourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicated_customer.test", "name", "test"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "replicated_customer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCustomerResourceConfig("test_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicated_customer.test", "name", "test_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomerResourceConfig(name string) string {
	return fmt.Sprintf(`
		resource "replicated_customer" "test" {
			name                    = %[1]q
			email                   = "test_resource@mm.mm"
			app_id                  = "test_app_id"
			channel_id              = "test_channel_id"
			expires_at              = "2025-01-30T15:04:05Z"
			is_kots_install_enabled = true

			entitlement_values = {
				testEntitlement	= "test_value"
			}
		}
	`, name)
}
