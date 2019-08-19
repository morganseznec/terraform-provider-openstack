package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/firewalls"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFWFirewallV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWFirewallV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: descriptions["tenant_id"],
			},

			"firewall_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"firewall_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceFWFirewallV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := firewalls.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("firewall_policy_id"); ok {
		listOpts.PolicyID = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOkExists("admin_state_up"); ok {
		listOpts.AdminStateUp = v.(bool)
	}

	if v, ok := d.GetOk("firewall_id"); ok {
		listOpts.ID = v.(string)
	}

	pages, err := firewalls.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_fw_firewall_v1: %s", err)
	}

	allFWFirewalls, err := firewalls.ExtractFirewalls(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract openstack_fw_firewall_v1: %s", err)
	}

	if len(allFWFirewalls) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allFWFirewalls) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	firewall := allFWFirewalls[0]

	log.Printf("[DEBUG] Retrieved openstack_fw_firewall_v1 %s: %#v", firewall.ID, firewall)
	d.SetId(firewall.ID)

	d.SetId(firewall.ID)
	d.Set("name", firewall.Name)
	d.Set("tenant_id", firewall.TenantID)
	d.Set("description", firewall.Description)
	d.Set("admin_state_up", firewall.AdminStateUp)
	d.Set("status", firewall.Status)
	d.Set("firewall_policy_id", firewall.PolicyID)
	d.Set("region", GetRegion(d, config))

	return nil
}
