package licensemanager_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func testAccGrantsDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	datasourceName := "data.aws_licensemanager_grants.test"
	licenseARN := os.Getenv(licenseARNKey)
	if licenseARN == "" {
		t.Skipf("Environment variable %s is not set to true", licenseARNKey)
	}
	principal := os.Getenv(principalKey)
	if principal == "" {
		t.Skipf("Environment variable %s is not set", principalKey)
	}
	homeRegion := os.Getenv(homeRegionKey)
	if homeRegion == "" {
		t.Skipf("Environment variable %s is not set", homeRegionKey)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantsDataSourceConfig_basic(licenseARN, rName, principal, homeRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "arns.0"),
				),
			},
		},
	})
}

func testAccGrantsDataSource_noMatch(t *testing.T) {
	datasourceName := "data.aws_licensemanager_grants.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGrantsDataSourceConfig_noMatch(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "arns.#", "0"),
				),
			},
		},
	})
}

func testAccGrantsDataSourceConfig_basic(licenseARN, rName, principal, homeRegion string) string {
	return acctest.ConfigCompose(
		acctest.ConfigRegionalProvider(homeRegion),
		acctest.ConfigAlternateAccountProvider(),
		fmt.Sprintf(`
data "aws_licensemanager_received_license" "test" {
  license_arn = %[1]q
}

locals {
  allowed_operations = [for i in data.aws_licensemanager_received_license.test.received_metadata[0].allowed_operations : i if i != "CreateGrant"]
}

resource "aws_licensemanager_grant" "test" {
  name               = %[2]q
  allowed_operations = local.allowed_operations
  license_arn        = data.aws_licensemanager_received_license.test.license_arn
  principal          = %[3]q
}

data "aws_licensemanager_grants" "test" {}
`, licenseARN, rName, principal),
	)
}

func testAccGrantsDataSourceConfig_noMatch() string {
	return `
data "aws_licensemanager_grants" "test" {
  filter {
    name = "LicenseIssuerName"
    values = [
      "No Match"
    ]
  }
}
`
}
