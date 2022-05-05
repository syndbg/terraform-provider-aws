package iam_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccIAMInstanceProfilesDataSource_basic(t *testing.T) {
	datasourceName := "data.aws_iam_instance_profiles.test"
	resourceName := "aws_iam_instance_profile.test"

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceProfilesDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "arns.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "paths.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "names.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceName, "arns.*", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "paths.*", resourceName, "path"),
					resource.TestCheckResourceAttrPair(datasourceName, "names.*", resourceName, "name"),
				),
			},
		},
	})
}

func testAccInstanceProfilesDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = %[1]q
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_instance_profile" "test" {
  name = "%[1]s-instance-profile"
  role = aws_iam_role.test.name
  path = "/testpath/"
}

data "aws_iam_instance_profiles" "test" {
  role_name  = aws_iam_role.test.name
  depends_on = [aws_iam_role.test, aws_iam_instance_profile.test]
}
`, rName)
}
