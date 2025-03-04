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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

const testAccClusternodegroup_authenticationvserver_binding_basic = `

resource "citrixadc_clusternodegroup_authenticationvserver_binding" "tf_clusternodegroup_authenticationvserver_binding" {
	name     = "my_authentication_group"
	vserver = "my_authentication_server"
  }
  
`

const testAccClusternodegroup_authenticationvserver_binding_basic_step2 = `
	# Keep the above bound resources without the actual binding to check proper deletion
`

func TestAccClusternodegroup_authenticationvserver_binding_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClusternodegroup_authenticationvserver_bindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccClusternodegroup_authenticationvserver_binding_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusternodegroup_authenticationvserver_bindingExist("citrixadc_clusternodegroup_authenticationvserver_binding.tf_clusternodegroup_authenticationvserver_binding", nil),
				),
			},
			resource.TestStep{
				Config: testAccClusternodegroup_authenticationvserver_binding_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusternodegroup_authenticationvserver_bindingNotExist("citrixadc_clusternodegroup_authenticationvserver_binding.tf_clusternodegroup_authenticationvserver_binding", "my_authentication_group,my_authentication_server"),
				),
			},
		},
	})
}

func testAccCheckClusternodegroup_authenticationvserver_bindingExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No clusternodegroup_authenticationvserver_binding id is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		client := testAccProvider.Meta().(*NetScalerNitroClient).client

		bindingId := rs.Primary.ID

		idSlice := strings.SplitN(bindingId, ",", 2)

		name := idSlice[0]
		vserver := idSlice[1]

		findParams := service.FindParams{
			ResourceType:             "clusternodegroup_authenticationvserver_binding",
			ResourceName:             name,
			ResourceMissingErrorCode: 258,
		}
		dataArr, err := client.FindResourceArrayWithParams(findParams)

		// Unexpected error
		if err != nil {
			return err
		}

		// Iterate through results to find the one with the matching vserver
		found := false
		for _, v := range dataArr {
			if v["vserver"].(string) == vserver {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("clusternodegroup_authenticationvserver_binding %s not found", n)
		}

		return nil
	}
}

func testAccCheckClusternodegroup_authenticationvserver_bindingNotExist(n string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*NetScalerNitroClient).client

		if !strings.Contains(id, ",") {
			return fmt.Errorf("Invalid id string %v. The id string must contain a comma.", id)
		}
		idSlice := strings.SplitN(id, ",", 2)

		name := idSlice[0]
		vserver := idSlice[1]

		findParams := service.FindParams{
			ResourceType:             "clusternodegroup_authenticationvserver_binding",
			ResourceName:             name,
			ResourceMissingErrorCode: 258,
		}
		dataArr, err := client.FindResourceArrayWithParams(findParams)

		// Unexpected error
		if err != nil {
			return err
		}

		// Iterate through results to hopefully not find the one with the matching vserver
		found := false
		for _, v := range dataArr {
			if v["vserver"].(string) == vserver {
				found = true
				break
			}
		}

		if found {
			return fmt.Errorf("clusternodegroup_authenticationvserver_binding %s was found, but it should have been destroyed", n)
		}

		return nil
	}
}

func testAccCheckClusternodegroup_authenticationvserver_bindingDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_clusternodegroup_authenticationvserver_binding" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(service.Clusternodegroup_authenticationvserver_binding.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clusternodegroup_authenticationvserver_binding %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
