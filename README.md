# terraform-provider-ipify
---

Provider to interact with the [ipify.org](https://ipify.org/) service.  Returns your IPv4 or IPv6 public IP in standard and CIDR notation.
Use the CIDR notation for securing cloud infrastructure with firewall rules when testing builds.

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