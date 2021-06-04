package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/atest"
)

func TestAccDataSourceAwsSnsTopic_basic(t *testing.T) {
	resourceName := "aws_sns_topic.test"
	datasourceName := "data.aws_sns_topic.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { atest.PreCheck(t) },
		ErrorCheck: atest.ErrorCheck(t, sns.EndpointsID),
		Providers:  atest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsSnsTopicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
				),
			},
		},
	})
}

func testAccDataSourceAwsSnsTopicConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "test" {
  name = %[1]q
}

data "aws_sns_topic" "test" {
  name = aws_sns_topic.test.name
}
`, rName)
}
