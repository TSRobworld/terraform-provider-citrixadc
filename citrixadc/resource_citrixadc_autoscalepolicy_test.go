/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package citrixadc

import (
	"fmt"
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const testAccAutoscalepolicy_basic = `


resource "citrixadc_autoscalepolicy" "tf_autoscalepolicy" {
	name         = "my_autoscaleprofile"
	rule         = "true"
	action       = "my_autoscaleaction"
  }
`
const testAccAutoscalepolicy_update = `


resource "citrixadc_autoscalepolicy" "tf_autoscalepolicy" {
	name         = "my_autoscaleprofile"
	rule         = "false"
	action       = "my_autoscaleaction"
  }
`
func TestAccAutoscalepolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAutoscalepolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAutoscalepolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalepolicyExist("citrixadc_autoscalepolicy.tf_autoscalepolicy", nil),
					resource.TestCheckResourceAttr("citrixadc_autoscalepolicy.tf_autoscalepolicy", "name", "my_autoscaleprofile"),
					resource.TestCheckResourceAttr("citrixadc_autoscalepolicy.tf_autoscalepolicy", "rule", "true"),
					resource.TestCheckResourceAttr("citrixadc_autoscalepolicy.tf_autoscalepolicy", "action", "my_autoscaleaction"),
				),
			},
			resource.TestStep{
				Config: testAccAutoscalepolicy_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalepolicyExist("citrixadc_autoscalepolicy.tf_autoscalepolicy", nil),
					resource.TestCheckResourceAttr("citrixadc_autoscalepolicy.tf_autoscalepolicy", "name", "my_autoscaleprofile"),
					resource.TestCheckResourceAttr("citrixadc_autoscalepolicy.tf_autoscalepolicy", "rule", "false"),
					resource.TestCheckResourceAttr("citrixadc_autoscalepolicy.tf_autoscalepolicy", "action", "my_autoscaleaction"),
				),
			},
		},
	})
}

func testAccCheckAutoscalepolicyExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No autoscalepolicy name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(service.Autoscalepolicy.Type(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("autoscalepolicy %s not found", n)
		}

		return nil
	}
}

func testAccCheckAutoscalepolicyDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_autoscalepolicy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(service.Autoscalepolicy.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("autoscalepolicy %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
