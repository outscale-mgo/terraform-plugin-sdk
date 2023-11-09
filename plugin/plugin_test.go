package plugin

import (
	"github.com/outscale-mgo/terraform-plugin-sdk/terraform"
)

func testProviderFixed(p terraform.ResourceProvider) ProviderFunc {
	return func() terraform.ResourceProvider {
		return p
	}
}
