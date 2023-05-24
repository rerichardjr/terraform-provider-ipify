terraform {
  required_providers {
    ipify = {
	version = "1.0.0"
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