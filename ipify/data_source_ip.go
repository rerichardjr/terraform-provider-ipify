package ipify

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ipifyURL = "https://api64.ipify.org/?format=json"

func dataSourceIP() *schema.Resource {
	return &schema.Resource{
		// Read: dataSourceIPRead,
		ReadContext: dataSourceIPRead,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
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

	d.SetId(message["ip"].(string))

	return diags
}
