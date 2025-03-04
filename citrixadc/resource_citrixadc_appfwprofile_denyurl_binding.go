package citrixadc

import (
	"net/url"

	"github.com/citrix/adc-nitro-go/resource/config/appfw"

	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
	"log"
	"strings"
)

func resourceCitrixAdcAppfwprofileDenyurlBinding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createAppfwprofileDenyurlBindingFunc,
		Read:          readAppfwprofileDenyurlBindingFunc,
		Delete:        deleteAppfwprofileDenyurlBindingFunc,
		Schema: map[string]*schema.Schema{
			"alertonly": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"denyurl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"isautodeployed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func createAppfwprofileDenyurlBindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In createAppfwprofileDenyurlBindingFunc")
	client := meta.(*NetScalerNitroClient).client

	profileName := d.Get("name")
	denyURL := d.Get("denyurl")

	// Use `,` as the separator since it is invalid character for adc entity strings
	bindingID := fmt.Sprintf("%s,%s", profileName, denyURL)

	appfwprofileDenyurlBinding := appfw.Appfwprofiledenyurlbinding{
		Alertonly:      d.Get("alertonly").(string),
		Comment:        d.Get("comment").(string),
		Denyurl:        d.Get("denyurl").(string),
		Isautodeployed: d.Get("isautodeployed").(string),
		Name:           d.Get("name").(string),
		State:          d.Get("state").(string),
	}

	err := client.UpdateUnnamedResource(service.Appfwprofile_denyurl_binding.Type(), &appfwprofileDenyurlBinding)
	if err != nil {
		return err
	}

	d.SetId(bindingID)

	err = readAppfwprofileDenyurlBindingFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this appfwprofileDenyurlBinding but we can't read it ?? %s", bindingID)
		return nil
	}
	return nil
}

func readAppfwprofileDenyurlBindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] citrixadc-provider:  In readAppfwprofileDenyurlBindingFunc")
	client := meta.(*NetScalerNitroClient).client

	bindingID := d.Id()
	idSlice := strings.SplitN(bindingID, ",", 2)

	if len(idSlice) < 2 {
		return fmt.Errorf("Cannot deduce appfwprofile and denyurl from ID string")
	}

	profileName := idSlice[0]
	denyURL := idSlice[1]

	log.Printf("[DEBUG] citrixadc-provider: Reading appfwprofileDenyurlBinding state %s", bindingID)

	findParams := service.FindParams{
		ResourceType: service.Appfwprofile_denyurl_binding.Type(),
		ResourceName: profileName,
	}
	findParams.FilterMap = make(map[string]string)
	findParams.FilterMap["denyurl"] = url.QueryEscape(denyURL)
	dataArr, err := client.FindResourceArrayWithParams(findParams)

	// Unexpected error
	if err != nil {
		log.Printf("[DEBUG] citrixadc-provider: Error during FindResourceArrayWithParams %s", err.Error())
		return err
	}

	// Resource is missing
	if len(dataArr) == 0 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams returned empty array")
		log.Printf("[WARN] citrixadc-provider: Clearing appfwprofile_denyurl_binding state %s", bindingID)
		d.SetId("")
		return nil
	}

	data := dataArr[0]

	d.Set("name", data["name"])
	d.Set("alertonly", data["alertonly"])
	d.Set("comment", data["comment"])
	d.Set("denyurl", data["denyurl"])
	d.Set("isautodeployed", data["isautodeployed"])
	d.Set("name", data["name"])
	d.Set("state", data["state"])

	return nil

}

func deleteAppfwprofileDenyurlBindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In deleteAppfwprofileDenyurlBindingFunc")
	client := meta.(*NetScalerNitroClient).client

	bindingID := d.Id()
	idSlice := strings.SplitN(bindingID, ",", 2)

	if len(idSlice) < 2 {
		return fmt.Errorf("Cannot deduce appfwprofile and denyurl from ID string")
	}

	profileName := idSlice[0]
	denyURL := idSlice[1]

	args := make([]string, 0)
	args = append(args, fmt.Sprintf("denyurl:%v", url.QueryEscape(denyURL)))

	err := client.DeleteResourceWithArgs(service.Appfwprofile_denyurl_binding.Type(), profileName, args)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
