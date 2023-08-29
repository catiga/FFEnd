package util

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

// var OssClient *oss.Client

// func OssInit() {
// 	OssClient, err := oss.New("ap-southeast-1", "LTAI5tQrAQHpYo9zfv9V6sVC", "LWI585O5wCabSen8v94f7wQJ0Btuyn") //新主账号
// 	if err != nil || OssClient == nil {
// 		fmt.Println("init OssClient error :", err)
// 		os.Exit(-1)
// 	}
// }

func GetSts() string {

	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	//       AI算命
	client, err := sts.NewClientWithAccessKey("ap-southeast-1", "LTAI5tN617yfYhiMhhgZeq5S", "9t4uTsivpsbTk3Buc3kFW8tSg92LTR")

	fmt.Println(client.GetConfig())

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	// request.Method = "GET"

	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	request.RoleArn = "acs:ram::5097164356872557:role/aliyunosstokengeneratorrole"
	request.RoleSessionName = "vegauser"
	request.DurationSeconds = requests.NewInteger(3600)

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	// response.
	fmt.Printf("response is %#v\n", response.GetHttpContentString())

	return response.GetHttpContentString()

}
