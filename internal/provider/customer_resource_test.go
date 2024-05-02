package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	rtypes "github.com/replicatedhq/replicated/pkg/types"
	"github.com/replicatedhq/replicated/pkg/util"
	"github.com/stretchr/testify/assert"
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

func TestGetCustomerResourceModelFromCustomer(t *testing.T) {
	expires, err := util.ParseTime("2025-01-30T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name                      string
		appID                     string
		customer                  *rtypes.Customer
		wantCustomerResourceModel CustomerResourceModel
	}{
		{
			name:  "customer with expires_at",
			appID: "123456789012",
			customer: &rtypes.Customer{
				ID:   "test_id",
				Name: "test_name",
				Channels: []rtypes.Channel{
					{
						ID: "test_channel_id",
					},
				},
				Expires: &util.Time{Time: expires},
			},
			wantCustomerResourceModel: CustomerResourceModel{
				AppId:     types.StringValue("123456789012"),
				Id:        types.StringValue("app/123456789012/customer/test_id"),
				Name:      types.StringValue("test_name"),
				ExpiresAt: types.StringValue("2025-01-30T15:04:05Z"),
			},
		},
		{
			name:  "customer without expires_at",
			appID: "123456789012",
			customer: &rtypes.Customer{
				ID: "test_id",
				Channels: []rtypes.Channel{
					{
						ID: "test_channel_id",
					},
				},
				Name: "test_name",
			},
			wantCustomerResourceModel: CustomerResourceModel{
				AppId: types.StringValue("123456789012"),
				Id:    types.StringValue("app/123456789012/customer/test_id"),
				Name:  types.StringValue("test_name"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCustomerResourceModel := getCustomerResourceModelFromCustomer(tt.appID, tt.customer)

			assert.Equal(t, tt.wantCustomerResourceModel.AppId, gotCustomerResourceModel.AppId)
			assert.Equal(t, tt.wantCustomerResourceModel.Id, gotCustomerResourceModel.Id)
			assert.Equal(t, tt.wantCustomerResourceModel.Name, gotCustomerResourceModel.Name)
		})
	}
}
