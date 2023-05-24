data "ipify_ip" "public" {}

output "public_ip" {
	description = "My Public IP"
	value = data.ipify_ip.public.ip_cidr
}