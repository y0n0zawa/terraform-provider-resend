package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainDataSource_basic(t *testing.T) {
	domainName := fmt.Sprintf("tf-test-%s.example.com", acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum))
	resourceName := "data.resend_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainDataSourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east-1"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

func testAccDomainDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "resend_domain" "test" {
  name = %[1]q
}

data "resend_domain" "test" {
  id = resend_domain.test.id
}
`, name)
}
