package tencentcloud

import (
	"context"
	"log"

	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ApmService struct {
	client *connectivity.TencentCloudClient
}

func (me *ApmService) DescribeApmById(ctx context.Context, instanceId string) (apmInstance *apm.ApmInstanceDetail, errRet error) {
	logId := getLogId(ctx)
	request := apm.NewDescribeApmInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseApmClient().DescribeApmInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Instances) < 1 {
		return
	}
	apmInstance = response.Response.Instances[0]
	return
}

func (me *ApmService) DescribeApmAgentById(ctx context.Context, instanceId string) (apmAgentInfo *apm.ApmAgentInfo, errRet error) {
	logId := getLogId(ctx)
	request := apm.NewDescribeApmAgentRequest()
	agentType := "go" //api有异常，先固定为go,不影响token获取
	request.InstanceId = &instanceId
	request.AgentType = &agentType

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseApmClient().DescribeApmAgent(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ApmAgent == nil {
		return
	}
	apmAgentInfo = response.Response.ApmAgent
	return
}
