# terraform-provider-ipify
---

Provider that interacts with [ipify.org](https://ipify.org/) and returns your IPv4 or IPv6 public IP in standard and CIDR notation.
Use the CIDR notation in firewall rules when you want to secure cloud infrastructure test builds.

## Using the provider

```hcl

terraform {
  required_providers {
    ipify = {
      source = "rerichardjr/ipify"
    }
  }
}

provider "ipify" {}

data "ipify_ip" "public" {}

output "public_ip" {
	description = "My Public IP"
	value = data.ipify_ip.public.ip_cidr
}
```