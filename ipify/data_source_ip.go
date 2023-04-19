package ipify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ipifyURL = "https://api64.ipify.org/?format=json"

func dataSourceIP() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for getting public IP address.",
		ReadContext: dataSourceIPRead,
		Schema: map[string]*schema.Schema{
			"ip": {
				Description: "The public IP address.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_cidr": {
				Description: "The public IP address in CIDR notation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest(http.MethodGet, ipifyURL, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	message := make(map[string]interface{}, 0)

	err = json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip", message["ip"]); err != nil {
		return diag.FromErr(err)
	}

	ipAddress := net.ParseIP(message["ip"].(string))

	if ipAddressType.To4() != nil {
		// IP is IPv4, set CIDR to /32
		if err := d.Set("ip_cidr", fmt.Sprintf("%s/32", ipAddress)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		// IP is IPv6, set CIDR to /128
		if err := d.Set("ip_cidr", fmt.Sprintf("%s/128", ipAddress)); err != nil {
			return diag.FromErr(err)
		}
	}
	
	d.SetId(message["ip"].(string))

	return diags
}
