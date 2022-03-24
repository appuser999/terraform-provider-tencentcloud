terraform {
  required_providers {
    tencentcloud = {
      source = "registry.terraform.io/tencentcloudstack/tencentcloud"
      version = ">= 1.0"
    }
  }
}

resource "tencentcloud_apm" "test-apm" {
  name          = "test-apm"
}

output "instance_id" {
  value = tencentcloud_apm.test-apm.id
}

data "tencentcloud_apm_agent" "test-apm-agent" {
  instance_id = tencentcloud_apm.test-apm.id
}

output "apm_agent_info" {
  value = data.tencentcloud_apm_agent.test-apm-agent
}