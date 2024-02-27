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
				Config: testAccCustomerResourceConfig("test3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicated_customer.test", "name", "test3"),
					// resource.TestCheckResourceAttrSet("replicated_customer.test", "version"),
					// resource.TestCheckResourceAttrSet("replicated_customer.test", "id"),
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
				Config: testAccCustomerResourceConfig("test1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicated_customer.test", "name", "test1"),
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
			email                   = "mcalzado+replicatedterraform@cartodb.com"
			app_id                  = "2WWnRhDha5d7j2deUt9Ejrx9igZ"
			channel_id              = "2aUOEOl1VgMlyma1xCr5KE933hO"
			expires_at              = "2024-01-30T15:04:05Z"
			is_kots_install_enabled = true

			entitlement_values = {
				cartoPlatformDefaultSA    = "abc"
				installerFeaturesEnabled  = "ac"
				selfHostedId              = "a"
				cartoFeaturesFlagSdkKey   = "a"
			}
		}
	`, name)
}
