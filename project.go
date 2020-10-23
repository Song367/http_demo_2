package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httppoint "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
	"http_demo_2/Services"
	"io"
	"net/url"
	"os"
	"time"
)

// 直连
func main2()  {
	target,_:=url.Parse("http://127.0.0.1:8080")
	//NewClient() go-kit/http引用的函数  最后两个func 一个是如何请求 ， 另一个是响应如何让处理
	client:=httppoint.NewClient("GET",target,Services.GetUserInfoRequest,Services.GetUserInfoResponse) //创建客户端
	getUserInfo:=client.Endpoint()      // 服务调用
	ctx:=context.Background()   //  创建空的上下文呢对象

	res,err:=getUserInfo(ctx,Services.UserRequest{UserId: 101})
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	UserRes:=res.(Services.UserResponse) //类型断言

	fmt.Println(UserRes)
}

// 通过consul来服务发现
func main(){
	{
		config := api.DefaultConfig()
		config.Address = "127.0.0.1:8500"
		apiClient, _ := api.NewClient(config)
		client:=consul.NewClient(apiClient)

		var logger log.Logger
		{
			logger=log.NewLogfmtLogger(os.Stdin)
		}
		{
			tags:=[]string{"primary"}
			Instance :=consul.NewInstancer(client,logger,"userService",tags,true)
			{
				factory:=func(Service_url string) (endpoint.Endpoint, io.Closer, error){
					target,_ := url.Parse("http://"+Service_url)
					return httppoint.NewClient("GET",target,Services.GetUserInfoRequest,Services.GetUserInfoResponse).Endpoint(),nil,nil
				}
				endpointer:=sd.NewEndpointer(Instance,factory,logger)
				Eps,err:=endpointer.Endpoints()
				if err!=nil{
					fmt.Println(err)
				}
				fmt.Println("endpoint 长度为",len(Eps))
				//myrb:=lb.NewRoundRobin(endpointer)  //  负载均衡调用服务 轮询方式
				myrb:=lb.NewRandom(endpointer,time.Now().UnixNano())  //  随机方式负载均衡
				for {

					getUserInfo,_:=myrb.Endpoint()  // 服务调用
					ctx:=context.Background()   //  创建空的上下文呢对象

					res,err:=getUserInfo(ctx,Services.UserRequest{UserId: 101})
					if err!=nil{
						fmt.Println(err)
						os.Exit(1)
					}
					UserRes:=res.(Services.UserResponse) //类型断言

					fmt.Println(UserRes.Result)
					time.Sleep(time.Second*3)
				}

			}
		}
	}
}
