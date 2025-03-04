package citrixadc

import (
	"github.com/citrix/adc-nitro-go/resource/config/network"
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
	"log"
	"net/url"
	"strings"
)

func resourceCitrixAdcNetbridge_nsip6_binding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createNetbridge_nsip6_bindingFunc,
		Read:          readNetbridge_nsip6_bindingFunc,
		Delete:        deleteNetbridge_nsip6_bindingFunc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"netmask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func createNetbridge_nsip6_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In createNetbridge_nsip6_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	name := d.Get("name").(string)
	ipaddress := d.Get("ipaddress").(string)
	bindingId := fmt.Sprintf("%s,%s", name, ipaddress)
	netbridge_nsip6_binding := network.Netbridgensip6binding{
		Ipaddress: d.Get("ipaddress").(string),
		Name:      d.Get("name").(string),
		Netmask:   d.Get("netmask").(string),
	}

	err := client.UpdateUnnamedResource(service.Netbridge_nsip6_binding.Type(), &netbridge_nsip6_binding)
	if err != nil {
		return err
	}

	d.SetId(bindingId)

	err = readNetbridge_nsip6_bindingFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this netbridge_nsip6_binding but we can't read it ?? %s", bindingId)
		return nil
	}
	return nil
}

func readNetbridge_nsip6_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] citrixadc-provider:  In readNetbridge_nsip6_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	name := idSlice[0]
	ipaddress := idSlice[1]
	log.Printf("IPADDRESS FROM SLICE:%s ", ipaddress)

	log.Printf("[DEBUG] citrixadc-provider: Reading netbridge_nsip6_binding state %s", bindingId)

	findParams := service.FindParams{
		ResourceType:             "netbridge_nsip6_binding",
		ResourceName:             name,
		ResourceMissingErrorCode: 258,
	}
	dataArr, err := client.FindResourceArrayWithParams(findParams)

	// Unexpected error
	if err != nil {
		log.Printf("[DEBUG] citrixadc-provider: Error during FindResourceArrayWithParams %s", err.Error())
		return err
	}

	// Resource is missing
	if len(dataArr) == 0 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams returned empty array")
		log.Printf("[WARN] citrixadc-provider: Clearing netbridge_nsip6_binding state %s", bindingId)
		d.SetId("")
		return nil
	}

	// Iterate through results to find the one with the right id
	foundIndex := -1
	for i, v := range dataArr {
		if v["ipaddress"].(string) == ipaddress {
			foundIndex = i
			break
		}
	}

	// Resource is missing
	if foundIndex == -1 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams ipaddress not found in array")
		log.Printf("[WARN] citrixadc-provider: Clearing netbridge_nsip6_binding state %s", bindingId)
		d.SetId("")
		return nil
	}
	// Fallthrough

	data := dataArr[foundIndex]

	d.Set("ipaddress", data["ipaddress"])
	d.Set("name", data["name"])
	d.Set("netmask", data["netmask"])

	return nil

}

func deleteNetbridge_nsip6_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In deleteNetbridge_nsip6_bindingFunc")
	client := meta.(*NetScalerNitroClient).client

	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	name := idSlice[0]
	ipaddress := idSlice[1]

	args := make([]string, 0)
	ipaddressEscaped := url.PathEscape(ipaddress)
	args = append(args, fmt.Sprintf("ipaddress:%s", ipaddressEscaped))
	//args = append(args, fmt.Sprintf("netmask:%s", d.Get("netmask").(string)))

	err := client.DeleteResourceWithArgs(service.Netbridge_nsip6_binding.Type(), name, args)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
