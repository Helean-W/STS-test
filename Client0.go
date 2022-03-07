package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/STS/snet"
)

//模拟客户端
func main() {
	fmt.Println("client start")

	time.Sleep(time.Second)
	//1直接连接，得到conn链接

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err:", err)
		return
	}

	//2链接调用Write写数据
	for {
		// _, err := conn.Write([]byte("Hello myserver..\n"))
		// if err != nil {
		// 	fmt.Println("write err:", err)
		// 	continue
		// }

		// buf := make([]byte, 512)

		// cnt, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("read err:", err)
		// 	continue
		// }

		// fmt.Printf("server return data: %s, cnt:%d", buf[:cnt], cnt)

		//发送封包的message消息, MsgID = 0
		dp := snet.NewDataPack()
		binaryMsg, err := dp.Pack(snet.NewMessage(0, []byte("sts client0 test msg....")))
		if err != nil {
			fmt.Println("pack error:", err)
			continue
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("server writer error:", err)
			continue
		}

		//发完服务器应该回复一个message数据
		//需要进行粘包处理

		//1先读取流中的head部分 得到ID 和 Len
		msgHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, msgHead); err != nil {
			fmt.Println("read head error:", err)
			break
		}
		msg, err := dp.Unpack(msgHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}

		//2再根据Len进行第二次读取，将data读出来
		if msg.GetMsgLen() > 0 {
			//msg里有数据，根据len再次读取
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("client read data error:", err)
				break
			}
			msg.SetData(data)

			fmt.Println("----->receive server message , msgID = ", msg.GetMsgId(), " , msgLen = ", msg.GetMsgLen(), " , data = ", string(msg.GetData()))
		}

		time.Sleep(time.Second)

	}
}
