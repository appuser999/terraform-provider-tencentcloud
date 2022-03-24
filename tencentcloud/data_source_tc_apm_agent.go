/*
Use this data source to query detailed information of audits.

Example Usage

```hcl
data "tencentcloud_audits" "audits" {
  name       = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
)

func dataSourceTencentCloudApmAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApmAgentRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the apm instance.",
			},
			"agent_download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the audits.",
			},
			"collector_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used to save results.",
			},
			"inner_collector_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used to save results.",
			},
			"private_link_collector_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used to save results.",
			},
			"public_collector_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used to save results.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudApmAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_apmagent.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	apmService := ApmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var agentInfo *apm.ApmAgentInfo
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		agentInfo, e = apmService.DescribeApmAgentById(ctx, instanceId)
		if e != nil {
			if strings.Contains(e.Error(), "TencentCloudSDKError") {
				return resource.NonRetryableError(e)
			}
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read apm instance agent failed, reason:%s\n ", logId, err.Error())
		return err
	}

	log.Printf("123: %v", *agentInfo.Token)

	_ = d.Set("agent_download_url", *agentInfo.AgentDownloadURL)
	_ = d.Set("collector_url", *agentInfo.CollectorURL)
	_ = d.Set("inner_collector_url", *agentInfo.InnerCollectorURL)
	_ = d.Set("private_link_collector_url", *agentInfo.PrivateLinkCollectorURL)
	_ = d.Set("public_collector_url", *agentInfo.PublicCollectorURL)
	_ = d.Set("token", *agentInfo.Token)

	return nil

}
