package main

import (
	"fmt"

	"github.com/STS/siface"
	"github.com/STS/snet"
)

// //Test PreHandle
// func (p *PingRouter) PreHandle(request siface.IRequest) {
// 	fmt.Println("Call Router PreHandle...")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
// 	if err != nil {
// 		fmt.Println("Call before ping error:", err)
// 	}
// }

// //Test Handle
// func (p *PingRouter) Handle(request siface.IRequest) {
// 	fmt.Println("Call Router Handle...")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
// 	if err != nil {
// 		fmt.Println("Call ping error:", err)
// 	}
// }

// //Test PostHandle
// func (p *PingRouter) PostHandle(request siface.IRequest) {
// 	fmt.Println("Call Router PostHandle...")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
// 	if err != nil {
// 		fmt.Println("Call after ping error:", err)
// 	}
// }

//ping test 自定义路由
type PingRouter struct {
	snet.BaseRouter
}

func (p *PingRouter) Handle(request siface.IRequest) {
	fmt.Println("Call Ping Router Handle...")
	//先读取客户端数据，再回写ping...ping...ping...
	fmt.Println("receive from client : msgID = ", request.GetMsgId(), ", data = ", string(request.GetMsgData()))
	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
	}
}

//hello test 自定义路由
type HelloStsRouter struct {
	snet.BaseRouter
}

func (h *HelloStsRouter) Handle(request siface.IRequest) {
	fmt.Println("Call Hello Router Handle...")
	//先读取客户端数据，再回写Hello...
	fmt.Println("receive from client : msgID = ", request.GetMsgId(), ", data = ", string(request.GetMsgData()))
	if err := request.GetConnection().SendMsg(201, []byte("Hello....")); err != nil {
		fmt.Println(err)
	}
}

//创建连接之后执行的Hook
func DoConnBegin(conn siface.IConnection) {
	fmt.Println("====>DoConnStart is Called...")
	if err := conn.SendMsg(202, []byte("DOCONNECTING BEGIN")); err != nil {
		fmt.Println(err)
	}

	//给当前连接设置一些属性
	conn.SetProperty("Name", "property测试。。。")
	conn.SetProperty("Gender", "male")
}

//连接断开之前执行的Hook
func DoConnEnd(conn siface.IConnection) {
	fmt.Println("====>DoConnEnd is Called...")
	fmt.Println("connection ID = ", conn.GetConnID(), "is offline...")

	//获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if gender, err := conn.GetProperty("Gender"); err == nil {
		fmt.Println("Gender = ", gender)
	}
}

func main() {
	//1 创建一个server句柄
	s := snet.NewServer("[zinx V0.1]")

	//2 注册连接Hook
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnEnd)

	//3 给当前框架添加自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloStsRouter{})

	//4 启动server
	s.Serve()
}
