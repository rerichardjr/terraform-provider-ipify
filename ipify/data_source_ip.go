package ipify

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

	publicIP := message["ip"].(string)

	// get the CIDR for the public IP, will be either 32 for IPv4 or 128 for IPv6
	cidr, err := getCIDR(publicIP)

	if err != nil {
		return diag.FromErr(err)
	}

	// set ip to the public IP	
	if err := d.Set("ip", publicIP); err != nil {
		return diag.FromErr(err)
	}	
	
	// set ip_cidr to the public IP plus the CIDR from getCIDR
	if err := d.Set("ip_cidr", fmt.Sprintf("%s/%s", publicIP, cidr)); err != nil {
		return diag.FromErr(err)
	}

	// set the id to the md5 hash of the public IP
	publicIPHash := md5.Sum([]byte(publicIP))
	d.SetId(hex.EncodeToString(publicIPHash[:]))

	return diags
}

func getCIDR(publicIP string) (cidr string, err error) {
	ipAddressType := net.ParseIP(publicIP)
	if ipAddressType == nil {
		return "", fmt.Errorf("the IP retrieved from %s doesn't appear to be valid", ipifyURL)
	}
	if ipAddressType.To4() != nil {
		return "32", nil
	}
	if ipAddressType.To16() != nil {
		return "128", nil
	}
	return "", fmt.Errorf("can't determine whether the IP retrieved from %s is IPv4 or IPv6", ipifyURL)
}
