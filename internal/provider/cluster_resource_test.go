package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("kind"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicated_cluster.test", "distribution", "kind"),
					resource.TestCheckResourceAttrSet("replicated_cluster.test", "version"),
					resource.TestCheckResourceAttrSet("replicated_cluster.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "replicated_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig("kind"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicated_cluster.test", "distribution", "kind"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(distribution string) string {
	return fmt.Sprintf(`
resource "replicated_cluster" "test" {
  distribution = %[1]q
}
`, distribution)
}
