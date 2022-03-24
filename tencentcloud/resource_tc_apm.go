/*
Provides a Load Balancer resource.

Example Usage

```hcl
resource "tencentcloud_apm" {
  name       = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudAPM() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPMCreate,
		Read:   resourceTencentCloudAPMRead,
		Update: resourceTencentCloudAPMUpdate,
		Delete: resourceTencentCloudAPMDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network type of the LB. Valid value: 'OPEN', 'INTERNAL'.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the LB. Valid value: 'CLASSIC', 'APPLICATION'.",
			},
			"trace_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The name of the LB.",
			},
			"span_daily_counters": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The VPC ID of the LB, unspecified or 0 stands for CVM basic network.",
			},
		},
	}
}

func resourceTencentCloudAPMCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apm.create")()

	logId := getLogId(contextNil)
	request := apm.NewCreateApmInstanceRequest()

	name := d.Get("name").(string)
	request.Name = helper.String(name)

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trace_duration"); ok {
		request.TraceDuration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("span_daily_counters"); ok {
		request.SpanDailyCounters = helper.Uint64(v.(uint64))
	}

	instanceId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check("create")
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseApmClient().CreateApmInstance(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response.Response.InstanceId == nil {
			err = fmt.Errorf("instance id is nil")
			return resource.NonRetryableError(err)
		}
		// requestId := *response.Response.RequestId
		instanceId = *response.Response.InstanceId

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apm instance failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(instanceId)

	return resourceTencentCloudAPMRead(d, meta)
}

func resourceTencentCloudAPMRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apm.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	apmService := ApmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *apm.ApmInstanceDetail
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := apmService.DescribeApmById(ctx, instanceId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read apm instance failed, reason:%s\n ", logId, err.Error())
		return err
	}

	_ = d.Set("name", instance.Name)
	_ = d.Set("description", instance.Description)
	_ = d.Set("trace_duration", instance.TraceDuration)
	_ = d.Set("span_daily_counters", instance.SpanDailyCounters)

	return nil
}

func resourceTencentCloudAPMUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apm.update")()

	return resourceTencentCloudAPMRead(d, meta)
}

func resourceTencentCloudAPMDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apm.delete")()

	return resourceTencentCloudAPMRead(d, meta)
}
